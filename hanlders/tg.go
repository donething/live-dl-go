package hanlders

import (
	"fmt"
	"github.com/donething/utils-go/dotg"
)

// FileSizeThreshold TG 视频大小限制（1.8GB）
const FileSizeThreshold = 1800 * 1024 * 1024

// TGHandler 发送到 TG
type TGHandler struct {
	TG     *dotg.TGBot // TGBot 的实例
	ChatID string      // 上传到的聊天频道
}

// Handle 发送**视频**到 TG
func (tg *TGHandler) Handle(task *TaskInfo) {
	_, err := tg.TG.SendVideo(tg.ChatID, task.Title, task.Path, task.FileSizeThreshold, "", false)

	if task.AfterHandle != nil {
		task.AfterHandle(task, err)
	}
}

// GenTGCaption 生成TG的标题Caption
//
// 参数为 主播名、所在平台名、日期、直播间标题，如：爱迟到的某、哔哩哔哩、20230415、进来看看
func GenTGCaption(name, plat, start string, roomTitle string) string {
	return dotg.LegalMk(fmt.Sprintf("#%s #%s %s _%s_", name, plat, start, roomTitle))
}
