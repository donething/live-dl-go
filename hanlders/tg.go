package hanlders

import (
	"fmt"
	"github.com/donething/utils-go/dotext"
	"github.com/donething/utils-go/dotg"
	"github.com/donething/utils-go/dovideo"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TGHandler 发送到 TG
type TGHandler struct {
	// TGBot 的实例
	TG *dotg.TGBot
	// 如果开启 TG 本地服务的端口
	LocalPort int
	// 上传到的聊天频道
	ChatID string
}

// Handle 发送**视频**到 TG
func (tg *TGHandler) Handle(info *InfoHandle) error {
	dst := info.Path
	// 不是 mp4 格式 的视频，才要转码为 mp4
	if strings.ToLower(filepath.Ext(info.Path)) != ".mp4" {
		dst = strings.TrimSuffix(info.Path, filepath.Ext(info.Path)) + ".mp4"
		err := dovideo.Convt(info.Path, dst)
		if err != nil {
			return fmt.Errorf("转码失败(%s)：\n%s", info.Path, err)
		}

		// 删除原视频。本来可以放在末尾的，但是占用磁盘空间，所以在转码成功后删除
		err = os.Remove(info.Path)
		if err != nil {
			return fmt.Errorf("删除原视频出错：%w", err)
		}
	}

	// 获取视频封面
	cover := strings.TrimSuffix(dst, filepath.Ext(dst)) + ".jpg"
	err := dovideo.GetFrame(dst, cover, "00:00:03", "320:320")
	if err != nil {
		return fmt.Errorf("获取视频封面出错(%s)：%s", dst, err)
	}

	// 上传到 TG

	// 数据
	// 直传
	vbs, err := os.Open(dst)
	if err != nil {
		return fmt.Errorf("读取视频文件出错：%w", err)
	}

	cbs, err := os.Open(cover)
	if err != nil {
		return fmt.Errorf("读取视频封面文件出错：%w", err)
	}

	w, h, err := dovideo.GetResolution(dst)
	if err != nil {
		return err
	}
	m := &dotg.InputMedia{
		MediaData: &dotg.MediaData{
			Type:              dotg.TypeVideo,
			Caption:           info.Title,
			ParseMode:         "MarkdownV2",
			Width:             w,
			Height:            h,
			SupportsStreaming: true,
		},
		Media:     vbs,
		Thumbnail: cbs,
	}
	_, err = tg.TG.SendMediaGroup(tg.ChatID, []*dotg.InputMedia{m})
	if err != nil {
		return fmt.Errorf("发送视频到TG出错：%w", err)
	}

	// 删除视频文件
	err = os.Remove(cover)
	if err != nil {
		return fmt.Errorf("删除视频封面出错：%w", err)
	}
	err = os.Remove(dst)
	if err != nil {
		return fmt.Errorf("删除转码后的视频出错：%w", err)
	}

	return nil
}

// GenTgCaption 生成TG的标题Caption
//
// 参数为 主播名、所在平台名、日期、直播间标题，如：爱迟到的某、哔哩哔哩、20230415、进来看看
func GenTgCaption(name, plat, roomTitle string) string {
	now := dotext.FormatDate(time.Now(), "20060102")
	return dotg.LegalMk(fmt.Sprintf("#%s #%s %s _%s_", name, plat, now, roomTitle))
}
