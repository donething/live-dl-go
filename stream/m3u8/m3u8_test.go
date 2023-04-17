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
		TG:        dotg.NewTGBot(os.Getenv("MY_TG_TOKEN")),
		LocalPort: 0,
		ChatID:    os.Getenv("MY_TG_CHAT_LIVE"),
	}
	title = dotg.EscapeMk(fmt.Sprintf("#测试 文件标题 %d", time.Now().UnixMilli()))
)

func TestStream_StartM3u8(t *testing.T) {
	u := "http://bjlive.szsbtech.com/record/vOowtobdl4dIdpe.m3u8?auth_key=1681753680-0-0-65d0eb722cfda380e42dae9cec4bbb1a"
	p := "D:/Tmp/live/zuji.ts"
	s := NewM3u8Stream(title, u, nil, p, 10*1024*1024, &tgHandler)

	err := s.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = <-s.GetChErr()
	if err != nil {
		t.Fatal(err)
	}

	restart := <-s.GetChRestart()
	t.Logf("重新下载直播流：%v", restart)
	if restart {
		TestStream_StartM3u8(t)
	}
}
