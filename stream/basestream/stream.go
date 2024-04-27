package basestream

import (
	"github.com/donething/live-dl-go/hanlders"
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
	*hanlders.TaskInfo

	StreamUrl string            // 视频流的地址
	Headers   map[string]string // 获取视频流时携带的请求头
	CurBytes  BytesType         // 已写入当前视频文件/文件夹的字节数，用于保证单个文件不超过指定的大小
	Stop      StopType          // 用于停止录制
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
