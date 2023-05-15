package flv

import (
	"fmt"
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/live-dl-go/stream/entity"
	"github.com/donething/live-dl-go/stream/files"
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
			Path:              path,
			FileSizeThreshold: fileSizeThreshold,
			Handler:           handler,
		},
	}
}

func (s *Stream) GetStream() *entity.Stream {
	return s.Stream
}

// Capture 录制 Flv 视频流
func (s *Stream) Capture() error {
	reader, err := s.CreateReader()
	if err != nil {
		return fmt.Errorf("创建 Flv 视频输入流出错：%w", err)
	}

	// 写入文件
	tFile := files.NewThresholdFile(reader, true, s.Path, s.FileSizeThreshold, s.Stream)

	err = tFile.StartSave()
	if err != nil {
		return fmt.Errorf("将 Flv 写入可限制大小的视频文件出错：%w", err)
	}

	return nil
}
