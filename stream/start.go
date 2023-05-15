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
	"github.com/donething/utils-go/domath"
	"github.com/donething/utils-go/dotext"
	"strings"
	"sync"
	"time"
)

const (
	// 获取主播信息失败的的最大次数
	maxFail = 3
)

// StartAnchor 开始录制直播流
//
// 参数为 正在录制表、直播流（Flv、M3u8）、主播信息、临时文件存储路径、单视频大小、视频处理器
//
// 当 stream 为 nil 时，将根据直播流地址自动生成
func StartAnchor(capturing *sync.Map, stream streamentity.IStream, anchor entity.Anchor, path string,
	fileSizeThreshold int64, handler hanlders.IHandler) error {
	// 开始录制该主播的时间
	start := dotext.FormatDate(time.Now(), "20060102")

	// 获取主播信息失败的次数
	fail := 0

	anchorSite, err := plats.GenAnchor(&anchor)
	if err != nil {
		return err
	}

	// 	换新文件保存视频，需要重新读取直播流的地址，以防旧的地址失效
LabelRetry:
	info, err := anchorSite.GetAnchorInfo()
	if err != nil {
		fail++

		// 重试
		if fail <= maxFail {
			logger.Warn.Printf("重试获取主播的信息(%+v)\n", anchor)
			time.Sleep(time.Duration(domath.RandInt(1, 3)) * time.Second)
			goto LabelRetry
		}

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
	// 存在表示正在录制且此次不用换新文件存储，不重复录制，返回
	if _, exists := capturing.Load(key); exists {
		logger.Info.Printf("😊【%s】正在录制(%+v)……\n", info.Name, anchor)
		return nil
	}

	// 需要开始录制

	// 生成标题
	// 平台对应的网站名
	title := hanlders.GenTgCaption(info.Name, anchorSite.GetPlatName(), start, info.Title)
	headers := anchorSite.GetStreamHeaders()

	// 如果没有指定直播流的类型，就自动匹配
	if stream == nil {
		if strings.Contains(strings.ToLower(info.StreamUrl), ".flv") {
			stream = flv.NewStream(title, info.StreamUrl, headers, path, fileSizeThreshold, handler)
		} else if strings.Contains(strings.ToLower(info.StreamUrl), ".m3u8") {
			stream = m3u8.NewStream(title, info.StreamUrl, headers, path, fileSizeThreshold, handler)
		} else {
			return fmt.Errorf("没有匹配到直播流的类型：%s", info.StreamUrl)
		}
	}

	// 开始录制直播流
	logger.Info.Printf("😙【%s】开始录制直播(%+v)\n", info.Name, anchor)

	// 记录正在录制的标识
	capturing.Store(key, stream)

	err = stream.Capture()
	if err != nil {
		capturing.Delete(key)
		return err
	}

	// 已下播，结束录制
	logger.Info.Printf("😶【%s】已中断直播(%+v)，停止录制\n", info.Name, anchor)
	capturing.Delete(key)

	return nil
}

// GenCapturingKey 正在录制的主播的键，避免重复录制，格式如 "<平台>_<主播ID>"，如 "bili_12345"
func GenCapturingKey(anchor *entity.Anchor) string {
	return fmt.Sprintf("%s_%s", anchor.Plat, anchor.ID)
}
