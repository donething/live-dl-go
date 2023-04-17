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

// NewFlvStream 创建 Stream 的实例
func NewFlvStream(title, streamUrl string, headers map[string]string, path string,
	fileSizeThreshold int, handler hanlders.IHandler) entity.IStream {
	return &Stream{
		Stream: &entity.Stream{
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

// Start 下载 flv 直播流
func (s *Stream) Start() error {
	err := s.PrepareCapture()
	if err != nil {
		return fmt.Errorf("准备录制flv流时出错：%w", err)
	}

	s.ChSegUrl <- s.LiveStreamUrl
	close(s.ChSegUrl)
	return nil
}

func (s *Stream) GetChErr() chan error {
	return s.ChErr
}

func (s *Stream) GetChRestart() chan bool {
	return s.ChRestart
}

func (s *Stream) Reset(title, streamUrl string, headers map[string]string, path string,
	fileSizeThreshold int, hanlder hanlders.IHandler) {
	s.ChErr = make(chan error)
	s.ChRestart = make(chan bool)
	s.ChSegUrl = make(chan string)

	s.Title = title
	s.LiveStreamUrl = streamUrl
	s.Headers = headers
	s.Path = path
	s.FileSizeThreshold = fileSizeThreshold
	s.Handler = hanlder
}
