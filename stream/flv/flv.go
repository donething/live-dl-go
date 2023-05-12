package flv

import (
	"fmt"
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/live-dl-go/stream/entity"
)

// Stream flv 直播流
type Stream struct {
	*entity.Stream
}

// NewStream 创建 Stream 的实例
func NewStream(title, streamUrl string, headers map[string]string, path string,
	fileSizeThreshold int64, handler hanlders.IHandler) entity.IStream {
	return &Stream{
		Stream: &entity.Stream{
			Title:             title,
			StreamUrl:         streamUrl,
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

// Start 下载 flv 直播流
func (s *Stream) Start() error {
	err := s.PrepareCapture()
	if err != nil {
		return fmt.Errorf("准备录制flv流时出错：%w", err)
	}

	s.ChSegUrl <- s.StreamUrl
	close(s.ChSegUrl)
	return nil
}

func (s *Stream) GetStream() *entity.Stream {
	return s.Stream
}

// Download 下载直播流、直链
func Download(title, streamUrl string, headers map[string]string,
	path string, handler hanlders.IHandler) error {
	s := NewStream(title, streamUrl, headers, path, 0, handler)
	err := s.Start()
	if err != nil {
		return err
	}

	err = <-s.GetStream().ChErr
	if err != nil {
		return err
	}

	return nil
}
