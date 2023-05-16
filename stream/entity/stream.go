package entity

import (
	"fmt"
	"github.com/donething/live-dl-go/comm"
	"github.com/donething/live-dl-go/hanlders"
	"io"
	"sync"
)

// BytesType 已写入的字节数
type BytesType struct {
	bytes int64
	mu    sync.Mutex
}

// StopType 停止录制视频流
type StopType struct {
	stop bool
	mu   sync.Mutex
}

// Stream 视频流
type Stream struct {
	// 标题。注意作为 TG 的 caption 时，需要转义
	Title string

	// 视频流的地址
	StreamUrl string

	// 请求头
	Headers map[string]string

	// 视频文件的保存路径
	Path string

	// 文件的最大字节数，为 0 表示无限制。上传 TG 建议设为 1.8GB: 1800*1024*1024
	FileSizeThreshold int64

	// 文件处理器
	Handler hanlders.IHandler

	// 	已写入当前视频文件/文件夹的字节数，用于保证单个文件不超过指定的大小
	CurBytes BytesType

	// 停止录制视频流
	Stop StopType
}

// CreateReader 创建输入流
func (s *Stream) CreateReader() (io.ReadCloser, error) {
	resp, err := comm.Client.Get(s.StreamUrl, s.Headers)
	if err != nil {
		return nil, fmt.Errorf("创建视频输入流出错。请求视频出错：%w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("创建视频输入流出错。读取视频的响应码：%s (URL: %s)", resp.Status, s.StreamUrl)
	}

	return resp.Body, nil
}

// GetBytes 获取当前视频文件中已写入的字节数
func (b *BytesType) GetBytes() int64 {
	var n int64

	b.mu.Lock()
	n = b.bytes
	b.mu.Unlock()

	return n
}

func (b *BytesType) AddBytes(n int64) {
	b.mu.Lock()
	b.bytes += n
	b.mu.Unlock()
}

func (b *BytesType) ResetBytes() {
	b.mu.Lock()
	b.bytes = 0
	b.mu.Unlock()
}

// SetStop 设置停止录制视频流
func (p *StopType) SetStop() {
	p.mu.Lock()
	p.stop = true
	p.mu.Unlock()
}

func (p *StopType) GetStop() bool {
	var stop bool

	p.mu.Lock()
	stop = p.stop
	p.mu.Unlock()

	return stop
}
