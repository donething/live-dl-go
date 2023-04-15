package stream

import (
	"dl-live-go/hanlders"
	"fmt"
)

// FlvStream flv 直播流
type FlvStream struct {
	*Stream
}

// NewFlvStream 创建 FlvStream 的实例
func NewFlvStream(title, streamUrl string, headers map[string]string, path string,
	fileSizeThreshold int, handler hanlders.IHandler) IStream {
	return &FlvStream{
		Stream: &Stream{
			Title:             title,
			LiveStreamUrl:     streamUrl,
			Headers:           headers,
			ChSegUrl:          make(chan string),
			Path:              path,
			FileSizeThreshold: fileSizeThreshold,
			Handler:           handler,
			ChErr:             make(chan error),
			ChRestart:         make(chan bool),
		},
	}
}

// StartAnchor 下载 flv 直播流
func (s *FlvStream) Start() error {
	err := s.PrepareCapture()
	if err != nil {
		return fmt.Errorf("准备录制flv流时出错：%w", err)
	}

	s.ChSegUrl <- s.LiveStreamUrl
	close(s.ChSegUrl)
	return nil
}

func (s *FlvStream) GetChErr() chan error {
	return s.ChErr
}

func (s *FlvStream) GetChRestart() chan bool {
	return s.ChRestart
}
