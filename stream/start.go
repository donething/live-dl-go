package stream

import (
	"fmt"
	"github.com/donething/live-dl-go/comm/logger"
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/live-dl-go/sites/entity"
	"github.com/donething/live-dl-go/sites/plats"
	streamentity "github.com/donething/live-dl-go/stream/entity"
	"github.com/donething/live-dl-go/stream/entity/capture_status"
	"github.com/donething/live-dl-go/stream/flv"
	"github.com/donething/live-dl-go/stream/m3u8"
	"github.com/donething/utils-go/dotext"
	"path/filepath"
	"strings"
	"time"
)

const (
	// 获取主播信息失败后，重试的次数
	maxRetry = 3
)

// StartAnchor 开始录制直播流
//
// 参数为：正在录制表、直播流（Flv、M3u8）、主播信息、临时工作路径、单视频大小、视频处理器
//
// 录制表 capturing 通过传递，方便在调用处获取录制状态
//
// 当 stream 为 nil 时，将根据直播流地址自动生成
func StartAnchor(capturing *capture_status.Capture[streamentity.IStream],
	stream streamentity.IStream,
	anchor entity.Anchor, workdir string, fileSizeThreshold int64, handler hanlders.IHandler) error {
	// 开始录制该主播的时间
	start := dotext.FormatDate(time.Now(), "20060102")

	anchorSite, err := plats.GenAnchor(&anchor)
	if err != nil {
		return err
	}

	// 	获取主播信息
	info, err := tryGetAnchorInfo(anchorSite, maxRetry)
	if err != nil {
		return err
	}

	// 读取录播状态的键
	key := capture_status.GenCapturingKey(&anchor)

	if !info.IsLive {
		logger.Info.Printf("😴【%s】没有在播(%+v)\n", info.Name, anchor)
		capturing.Del(key)
		return nil
	}

	// 判断此次是否需要录制视频
	// 存在表示正在录制且此次不用换新文件存储，不重复录制，返回
	if s, exists := capturing.Get(key); exists {
		bytes := dotext.BytesHumanReadable(s.GetStream().CurBytes.GetBytes())
		logger.Info.Printf("😊【%s】正在录制(%+v)，当前文件已写入 %s/%s\n", info.Name, anchor,
			bytes, dotext.BytesHumanReadable(fileSizeThreshold))
		return nil
	}

	// 需要开始录制

	// 生成标题
	// 平台对应的网站名
	title := hanlders.GenTgCaption(info.Name, anchorSite.GetPlatName(), start, info.Title)
	headers := anchorSite.GetStreamHeaders()

	// 如果没有指定直播流的类型，就自动匹配
	if stream == nil {
		name := fmt.Sprintf("%s_%s", anchor.Plat, anchor.ID)
		if strings.Contains(strings.ToLower(info.StreamUrl), ".flv") {
			// 保存依然为 flv，只是发送到 TG 前转为 mp4
			path := filepath.Join(workdir, name+".flv")
			stream = flv.NewStream(title, info.StreamUrl, headers, path, fileSizeThreshold, handler)
		} else if strings.Contains(strings.ToLower(info.StreamUrl), ".m3u8") {
			// m3u8 合并片段时就转为 mp4
			path := filepath.Join(workdir, name+".mp4")
			stream = m3u8.NewStream(title, info.StreamUrl, headers, path, fileSizeThreshold, handler)
		} else {
			return fmt.Errorf("没有匹配到直播流的类型：%s", info.StreamUrl)
		}
	}
	// 记录正在录制的标识
	capturing.Set(key, stream)

	// 开始录制直播流
	logger.Info.Printf("😙【%s】开始录制直播(%+v)\n", info.Name, anchor)

	err = stream.Capture()
	// 当录制出错时，要判断出错情况：在获取直播流出错时，先判断主播此时是否在播，主播且出错才是真正的录制错误
	if err != nil {
		infoCheck, err := tryGetAnchorInfo(anchorSite, maxRetry)
		if err != nil {
			return err
		}

		if infoCheck.IsLive {
			return err
		}
	}

	// 已下播，结束录制
	logger.Info.Printf("😶【%s】已中断直播(%+v)，停止录制\n", info.Name, anchor)
	capturing.Del(key)

	return nil
}

// 获取主播信息，可指定失败后的重试次数
func tryGetAnchorInfo(anchorSite entity.IAnchor, retry int) (*entity.AnchorInfo, error) {
	fail := 0
	var info *entity.AnchorInfo
	var err error

	for {
		info, err = anchorSite.GetAnchorInfo()
		if err != nil {
			// 重试
			if fail < retry {
				fail++
				time.Sleep(1 * time.Second)
				continue
			}

			return nil, err
		}

		// 获取成功
		return info, nil
	}
}
