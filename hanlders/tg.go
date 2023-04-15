package hanlders

import (
	"fmt"
	"github.com/donething/utils-go/dotext"
	"github.com/donething/utils-go/dotg"
	"os"
	"os/exec"
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

// Handle 发送到 TG
func (tg *TGHandler) Handle(info *InfoHandle) error {
	// 转码
	dst := info.Path + ".mp4"
	cover := info.Path + ".jpg"

	cmd := exec.Command("ffmpeg", "-hide_banner", "-i", info.Path, "-c", "copy", dst)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("转码失败(%s)：\n%s: %w", info.Path, string(output), err)
	}

	// 获取视频封面
	// 可以省略压缩图片："-vf", "scale=512:512:force_original_aspect_ratio=decrease"
	cmd = exec.Command("ffmpeg", "-hide_banner", "-i", info.Path, "-ss",
		"00:00:03", "-frames:v", "1", cover)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("获取视频封面出错：%s: %w", string(output), err)
	}

	// 删除原视频。本来可以放在末尾的，但是占用磁盘空间，所以在转码成功后删除
	err = os.Remove(info.Path)
	if err != nil {
		return fmt.Errorf("删除原视频出错：%w", err)
	}

	// 上传到 TG
	// name := filepath.Base(dst)
	// // 文件标题
	// now := dotext.FormatDate(time.Now(), "20060102")
	// caption := dotg.LegalMk(fmt.Sprintf("#%s #%s %s _%s_", info.Name, info.Plat, now, info.Title))
	// // 标题
	// open := fmt.Sprintf("[直播间](%s)", info.WebUrl)
	// caption = fmt.Sprintf("%s %s", caption, open)

	// 数据
	var vData interface{}
	var cData interface{}
	if tg.LocalPort != 0 {
		// 通过 TG 本地服务上传
		vData = fmt.Sprintf("file://%s", dst)
		cData = fmt.Sprintf("file://%s", cover)
	} else {
		// 直传
		bs, err := os.ReadFile(dst)
		if err != nil {
			return fmt.Errorf("读取视频文件出错：%w", err)
		}
		vData = bs

		bs, err = os.ReadFile(cover)
		if err != nil {
			return fmt.Errorf("读取视频封面文件出错：%w", err)
		}
		cData = bs
	}

	m := &dotg.InputMedia{
		Type:              dotg.TypeVideo,
		Media:             vData,
		Thumbnail:         cData,
		Caption:           info.Title,
		ParseMode:         "MarkdownV2",
		SupportsStreaming: true,
	}
	_, err = tg.TG.SendMediaGroup(tg.ChatID, []*dotg.InputMedia{m})
	if err != nil {
		return fmt.Errorf("发送视频出错到TG：%w", err)
	}

	err = os.Remove(cover)
	if err != nil {
		return fmt.Errorf("删除转码视频出错：%w", err)
	}

	err = os.Remove(dst)
	return err
}

// GenTgCaption 生成TG的标题Caption
//
// 参数为 主播名、所在平台名、日期、直播间标题，如：爱迟到的某、哔哩哔哩、20230415、进来看看
func GenTgCaption(name, plat, roomTitle string) string {
	now := dotext.FormatDate(time.Now(), "20060102")
	return dotg.LegalMk(fmt.Sprintf("#%s #%s %s _%s_", name, plat, now, roomTitle))
}
