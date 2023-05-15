package flv

import (
	"fmt"
	"github.com/donething/live-dl-go/comm"
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
	headers := map[string]string{
		// referer 必不可少
		"referer":    "https://live.bilibili.com/",
		"user-agent": comm.UAWin,
	}

	s := NewStream("哔哩616",
		"https://cn-gddg-cu-01-04.bilivideo.com/live-bvc/398404/live_50329118_9516950_bluray.flv?expires=1684107327&pt=web&deadline=1684107327&len=0&oi=1885426481&platform=web&qn=10000&trid=1000cde80f162c2e4e29b4bb751170dfcc83&uipk=100&uipv=100&nbs=1&uparams=cdn,deadline,len,oi,platform,qn,trid,uipk,uipv,nbs&cdn=cn-gotcha01&upsig=1b089e4b0d8690801ab3c02a8fb52c16&sk=1304f646dfeb4df8b6e7ff33c167d3ad7ba0ab83aecab529c31066fdbad2ad21&p2p_type=0&sl=10&free_type=0&mid=0&sid=cn-gddg-cu-01-04&chash=0&sche=ban&score=18&pp=rtmp&source=one&trace=8a0&site=5e68e623b6df60387d791cefa4ef1a3b&order=1", headers, "D:/Tmp/live/bili_616.flv", 5*1024*1024, &tgHandler)
	err := s.Capture()
	if err != nil {
		t.Fatal(err)
	}
}
