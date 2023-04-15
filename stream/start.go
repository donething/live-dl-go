package stream

import (
	"dl-live-go/hanlders"
	_ "dl-live-go/sites"
	"dl-live-go/sites/plats"
	"fmt"
)

// NewStreamType 函数类型，用于创建`Stream`
type NewStreamType func(title, streamUrl string, headers map[string]string, path string,
	fileSizeThreshold int, handler hanlders.IHandler) IStream

// StartAnchor 开始录制直播流
//
// 参数为 主播信息、临时文件存储路径、单视频大小、视频处理器
func StartAnchor(New NewStreamType, anchor *plats.Anchor, path string,
	fileSizeThreshold int, handler hanlders.IHandler) error {
	plat, ok := plats.Plats[anchor.Plat]
	if !ok {
		return fmt.Errorf("未知的平台'%s'", anchor.Plat)
	}

	info, err := plat.GetAnchorInfo(anchor.ID)
	if err != nil {
		return err
	}
	if !info.IsLive {
		return fmt.Errorf("用户不在播，无法开始录制")
	}

	// 生成标题
	// 平台对应的网站名
	site := plats.Sites[anchor.Plat]
	title := hanlders.GenTgCaption(info.Name, site, info.Title)
	headers := plats.Headers[anchor.Plat]
	// 创建实例
	stream := New(title, info.StreamUrl, headers, path, fileSizeThreshold, handler)

	// 开始录制直播流
	err = stream.Start()
	if err != nil {
		return err
	}

	// 等待下载阶段的错误
	err = <-stream.GetChErr()
	if err != nil {
		return err
	}

	// 需要用新的文件存储视频
	restart := <-stream.GetChRestart()
	if restart {
		return StartAnchor(New, anchor, path, fileSizeThreshold, handler)
	}

	return nil
}

// StartFlvAnchor 开始录制 flv 直播流
//
// 参数为 主播信息、临时文件存储路径、单视频大小、视频处理器
func StartFlvAnchor(anchor *plats.Anchor, path string, fileSizeThreshold int,
	handler hanlders.IHandler) error {
	return StartAnchor(NewFlvStream, anchor, path, fileSizeThreshold, handler)
}

// StartM3u8Anchor 开始录制 m3u8 直播流
//
// 参数为 主播信息、临时文件存储路径、单视频大小、视频处理器
func StartM3u8Anchor(anchor *plats.Anchor, path string, fileSizeThreshold int,
	handler hanlders.IHandler) error {
	return StartAnchor(NewM3u8Stream, anchor, path, fileSizeThreshold, handler)
}
