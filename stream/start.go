package stream

import (
	"fmt"
	"live-dl-go/comm/logger"
	"live-dl-go/hanlders"
	_ "live-dl-go/sites"
	"live-dl-go/sites/plats"
	"sync"
)

// NewStreamType 函数类型，用于创建`Stream`
type NewStreamType func(title, streamUrl string, headers map[string]string, path string,
	fileSizeThreshold int, handler hanlders.IHandler) IStream

// 开始录制直播流
//
// 参数为 主播信息、临时文件存储路径、单视频大小、视频处理器
func startAnchor(capturing *sync.Map, stream IStream, anchor *plats.Anchor, path string,
	fileSizeThreshold int, handler hanlders.IHandler) error {
	// 此次是否是换新文件保存视频
	// 用于当正在录播且isNewFile为真时，不退出
	var isNewFile = false

	plat, ok := plats.Plats[anchor.Plat]
	if !ok {
		return fmt.Errorf("未知的平台'%s'", anchor.Plat)
	}

	// 	换新文件保存视频，需要重新读取直播流的地址，以防旧的地址失效
LabelNewFile:
	info, err := plat.GetAnchorInfo(anchor.ID)
	if err != nil {
		return err
	}
	// 是否正在录播的键
	key := GenCapturingKey(anchor.Plat, anchor.ID)

	if !info.IsLive {
		logger.Info.Printf("主播未在播【%s】(%+v)\n", info.Name, *anchor)
		capturing.Delete(key)
		return nil
	}

	// 判断此次是否需要录制视频
	// 存在表示正在录制，不重复录制，返回
	if _, exists := capturing.Load(key); !isNewFile && exists {
		logger.Info.Printf("该直播间正在录制【%s】(%+v)\n", info.Name, *anchor)
		return nil
	}

	// 需要开始录制

	// 生成标题
	// 平台对应的网站名
	site := plats.Sites[anchor.Plat]
	title := hanlders.GenTgCaption(info.Name, site, info.Title)
	headers := plats.Headers[anchor.Plat]
	// 设置流的信息
	stream.Reset(title, info.StreamUrl, headers, path, fileSizeThreshold, handler)

	// 开始录制直播流
	logger.Info.Printf("开始录制直播间【%s】(%+v)\n", info.Name, info)
	err = stream.Start()
	if err != nil {
		return err
	}

	// 记录正在录制的标识
	capturing.Store(key, true)

	// 等待下载阶段的错误
	err = <-stream.GetChErr()
	if err != nil {
		capturing.Delete(key)
		return err
	}

	// 需要用新的文件存储视频
	restart := <-stream.GetChRestart()
	if restart {
		isNewFile = true
		goto LabelNewFile
	}

	// 已下播，结束录制
	logger.Info.Printf("直播间已中断直播【%s】(%+v)，停止录制\n", info.Name, anchor)
	capturing.Delete(key)

	return nil
}

// StartFlvAnchor 开始录制 flv 直播流
//
// 参数为 正在录制表、主播信息、临时文件存储路径（不需担心重名）、单视频大小、视频处理器
func StartFlvAnchor(capturing *sync.Map, anchor *plats.Anchor, path string, fileSizeThreshold int,
	handler hanlders.IHandler) error {
	s := &FlvStream{Stream: &Stream{
		Path:              path,
		FileSizeThreshold: fileSizeThreshold,
		Handler:           handler,
	}}

	return startAnchor(capturing, s, anchor, path, fileSizeThreshold, handler)
}

// StartM3u8Anchor 开始录制 m3u8 直播流
//
// 参数为 正在录制表、主播信息、临时文件存储路径（不需担心重名）、单视频大小、视频处理器
func StartM3u8Anchor(capturing *sync.Map, anchor *plats.Anchor, path string, fileSizeThreshold int,
	handler hanlders.IHandler) error {
	s := &M3u8Stream{Stream: &Stream{
		Path:              path,
		FileSizeThreshold: fileSizeThreshold,
		Handler:           handler,
	}}

	return startAnchor(capturing, s, anchor, path, fileSizeThreshold, handler)
}

// GenCapturingKey 正在录制的主播的键，避免重复录制，格式如 "<平台>_<主播ID>"，如 "bili_12345"
func GenCapturingKey(plat, id string) string {
	return fmt.Sprintf("%s_%s", plat, id)
}
