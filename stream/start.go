package stream

import (
	"fmt"
	"github.com/donething/live-dl-go/comm/logger"
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/live-dl-go/sites/entity"
	"github.com/donething/live-dl-go/sites/plats"
	streamentity "github.com/donething/live-dl-go/stream/entity"
	"github.com/donething/live-dl-go/stream/flv"
	"github.com/donething/live-dl-go/stream/m3u8"
	"strings"
	"sync"
)

// NewStreamType 函数类型，用于创建`Stream`
type NewStreamType func(title, streamUrl string, headers map[string]string, path string,
	fileSizeThreshold int, handler hanlders.IHandler) streamentity.IStream

// StartAnchor 开始录制直播流
//
// 参数为 正在录制表、直播流（Flv、M3u8）、主播信息、临时文件存储路径、单视频大小、视频处理器
//
// 当 stream 为 nil 时，将根据直播流地址自动生成
func StartAnchor(capturing *sync.Map, stream streamentity.IStream, anchor entity.Anchor, path string,
	fileSizeThreshold int, handler hanlders.IHandler) error {
	// 此次是否是换新文件保存视频
	// 用于当正在录播且isNewFile为真时，不退出
	var isNewFile = false

	anchorSite, err := plats.GenAnchor(&anchor)
	if err != nil {
		return err
	}

	// 	换新文件保存视频，需要重新读取直播流的地址，以防旧的地址失效
LabelNewFile:
	info, err := anchorSite.GetAnchorInfo()
	if err != nil {
		return err
	}
	// 是否正在录播的键
	key := GenCapturingKey(&anchor)

	if !info.IsLive {
		logger.Info.Printf("😴【%s】没有在播(%+v)\n", info.Name, anchor)
		capturing.Delete(key)
		return nil
	}

	// 判断此次是否需要录制视频
	// 存在表示正在录制，不重复录制，返回
	if _, exists := capturing.Load(key); !isNewFile && exists {
		logger.Info.Printf("😊【%s】正在录制…(%+v)\n", info.Name, anchor)
		return nil
	}

	// 需要开始录制

	// 生成标题
	// 平台对应的网站名
	title := hanlders.GenTgCaption(info.Name, anchorSite.GetPlatName(), info.Title)
	headers := anchorSite.GetStreamHeaders()

	// 如果没有指定直播流的类型，就自动匹配
	if stream == nil {
		if strings.Contains(strings.ToLower(info.StreamUrl), ".flv") {
			stream = &flv.Stream{Stream: &streamentity.Stream{}}
		} else if strings.Contains(strings.ToLower(info.StreamUrl), ".m3u8") {
			stream = &m3u8.Stream{Stream: &streamentity.Stream{}}
		} else {
			return fmt.Errorf("没有匹配到直播流的类型：%s", info.StreamUrl)
		}
	}

	// 设置流的信息
	stream.Reset(title, info.StreamUrl, headers, path, fileSizeThreshold, handler)

	// 开始录制直播流
	logger.Info.Printf("😙开始录制直播间【%s】(%+v)\n", info.Name, anchor)
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
	logger.Info.Printf("😶直播间已中断直播【%s】(%+v)，停止录制\n", info.Name, anchor)
	capturing.Delete(key)

	return nil
}

// StartFlvAnchor 开始录制 flv 直播流
//
// 参数为 正在录制表、主播信息、临时文件存储路径（不需担心重名）、单视频大小、视频处理器
func StartFlvAnchor(capturing *sync.Map, anchor entity.Anchor, path string, fileSizeThreshold int,
	handler hanlders.IHandler) error {
	s := &flv.Stream{Stream: &streamentity.Stream{}}

	return StartAnchor(capturing, s, anchor, path, fileSizeThreshold, handler)
}

// StartM3u8Anchor 开始录制 m3u8 直播流
//
// 参数为 正在录制表、主播信息、临时文件存储路径（不需担心重名）、单视频大小、视频处理器
//
// 下载m3u8视频（非直播）时，可下载到单个文件中，不能分文件保存，因为会重读m3u8文件，也就会重头开始下载
func StartM3u8Anchor(capturing *sync.Map, anchor entity.Anchor, path string, fileSizeThreshold int,
	handler hanlders.IHandler) error {
	s := &m3u8.Stream{Stream: &streamentity.Stream{}}

	return StartAnchor(capturing, s, anchor, path, fileSizeThreshold, handler)
}

// GenCapturingKey 正在录制的主播的键，避免重复录制，格式如 "<平台>_<主播ID>"，如 "bili_12345"
func GenCapturingKey(anchor *entity.Anchor) string {
	return fmt.Sprintf("%s_%s", anchor.Plat, anchor.ID)
}
