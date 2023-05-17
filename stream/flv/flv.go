package flv

import (
	"fmt"
	"github.com/donething/live-dl-go/comm"
	"github.com/donething/live-dl-go/hanlders"
	entity2 "github.com/donething/live-dl-go/sites/entity"
	"github.com/donething/live-dl-go/stream/entity"
	"io"
)

// CreateReaderFun 创建 flv 视频输入流的函数类型
type CreateReaderFun func(anchorSite entity2.IAnchor) (io.ReadCloser, error)

// Stream flv 直播流
type Stream struct {
	*entity.Stream
	anchor entity2.IAnchor
}

// NewStream 创建 Stream 的实例
//
// 参数 path 视频的保存路径，以 ".flv" 结尾
func NewStream(title, path string, fileSizeThreshold int64,
	handler hanlders.IHandler, anchor entity2.IAnchor) entity.IStream {
	return &Stream{
		Stream: &entity.Stream{
			Title:             title,
			Path:              path,
			FileSizeThreshold: fileSizeThreshold,
			Handler:           handler,
		},
		anchor: anchor,
	}
}

func (s *Stream) GetStream() *entity.Stream {
	return s.Stream
}

// Capture 录制 Flv 视频流
func (s *Stream) Capture() error {
	// 写入文件
	tFile := NewThresholdFile(s.Path, s.FileSizeThreshold, s)

	err := tFile.StartSave()
	if err != nil {
		return fmt.Errorf("将 Flv 视频流写入可限制大小的视频文件出错：%w", err)
	}

	return nil
}

// CreateReader 创建 flv 视频输入流
func (s *Stream) CreateReader() (io.ReadCloser, error) {
	info, err := entity2.TryGetAnchorInfo(s.anchor, entity2.MaxRetry)
	if err != nil {
		return nil, fmt.Errorf("创建 Flv Reader 出错：获取主播信息出错：%w", err)
	}

	resp, err := comm.Client.Get(info.StreamUrl, s.anchor.GetStreamHeaders())
	if err != nil {
		return nil, fmt.Errorf("创建 Flv Reader 出错：：请求视频出错：%w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("创建 Flv Reader 出错：：读取视频的响应码：%s (URL: %s)",
			resp.Status, info.StreamUrl)
	}

	return resp.Body, nil
}
