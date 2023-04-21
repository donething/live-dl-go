package hanlders

import (
	"fmt"
	"github.com/donething/utils-go/dotext"
	"github.com/donething/utils-go/dotg"
	"os"
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
	media, dst, thumb, err := dotg.GenTgMedia(info.Path, info.Title)
	_, err = tg.TG.SendMediaGroup(tg.ChatID, []*dotg.InputMedia{media})
	if err != nil {
		return fmt.Errorf("发送视频到TG出错：%w", err)
	}

	// 删除视频文件
	err = os.Remove(thumb)
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
