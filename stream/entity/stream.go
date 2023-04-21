package entity

import (
	"github.com/donething/live-dl-go/hanlders"
)

// Stream 直播流
type Stream struct {
	// 标题。注意作为TG的caption时，需要转义
	Title string
	// 直播流的地址
	StreamUrl string
	// 请求头
	Headers map[string]string
	// 用于发送视频（切片）的下载URL
	// m3u8 直播流需要经常更新m3u8文件获取新的切片，需要通过channel传递URL。缓冲池可以设为 5
	// flv 直播流可以不用设置缓冲池
	ChSegUrl chan string
	// 视频文件的保存路径
	Path string
	// 文件的最大字节数，为 0 表示无限制。建议 1GB: 1024*1024*1024
	FileSizeThreshold int64
	// 文件处理器
	Handler hanlders.IHandler
	// 因为下载视频到处理视频，要经过多个 goroutine，用 channel 传递错误信息
	ChErr chan error
	// 每保存一个视频文件就重新开始获取视频流。这样避免手动为视频添加头信息
	// 需要手动实现重新开始下载，参考 `StartAnchor`函数
	ChRestart chan bool
}

// Reset 重置信息。传递需要设置或修改的参数
func (s *Stream) Reset(title, streamUrl string, headers map[string]string, path string,
	fileSizeThreshold int64, hanlder hanlders.IHandler) {
	s.ChErr = make(chan error)
	s.ChRestart = make(chan bool)
	s.ChSegUrl = make(chan string)

	s.Title = title
	s.StreamUrl = streamUrl
	s.Headers = headers
	s.Path = path
	s.FileSizeThreshold = fileSizeThreshold
	s.Handler = hanlder
}
