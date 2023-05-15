package m3u8

import (
	"fmt"
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/utils-go/dotg"
	"os"
	"testing"
	"time"
)

var (
	tgHandler = hanlders.TGHandler{
		TG:     dotg.NewTGBot(os.Getenv("MY_TG_TOKEN")),
		ChatID: os.Getenv("MY_TG_CHAT_LIVE"),
	}
	title = dotg.EscapeMk(fmt.Sprintf("#测试 文件标题 %d", time.Now().UnixMilli()))
)

func TestStream_Capture(t *testing.T) {
	s := NewStream("足迹15722883",
		"http://bjlive.szsbtech.com/record/Nx8t2oAALlhNjN.m3u8?auth_key=1684146012-0-0-16b677731be5aed037a429763c8ea515",
		nil,
		"D:/Tmp/live/zuji_15722883.ts",
		5*1024*1024, &tgHandler)
	err := s.Capture()
	if err != nil {
		t.Fatal(err)
	}
}
