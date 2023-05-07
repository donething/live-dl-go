package hanlders

import (
	"fmt"
	"github.com/donething/utils-go/dotg"
)

// TGHandler 发送到 TG
type TGHandler struct {
	// TGBot 的实例
	TG *dotg.TGBot
	// 上传到的聊天频道
	ChatID string
}

// Handle 发送**视频**到 TG
func (tg *TGHandler) Handle(info *InfoHandle) error {
	_, err := tg.TG.SendVideo(tg.ChatID, info.Title, info.Path, info.FileSizeThreshold, "", info.Reserve)
	return err
}

// GenTgCaption 生成TG的标题Caption
//
// 参数为 主播名、所在平台名、日期、直播间标题，如：爱迟到的某、哔哩哔哩、20230415、进来看看
func GenTgCaption(name, plat, start string, roomTitle string) string {
	return dotg.LegalMk(fmt.Sprintf("#%s #%s %s _%s_", name, plat, start, roomTitle))
}
