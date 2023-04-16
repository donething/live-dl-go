package stream

import (
	"fmt"
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/live-dl-go/sites/plats"
	"github.com/donething/utils-go/dotg"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	capturing = &sync.Map{}

	tgHandler = hanlders.TGHandler{
		TG:        dotg.NewTGBot(os.Getenv("MY_TG_TOKEN")),
		LocalPort: 0,
		ChatID:    os.Getenv("MY_TG_CHAT_LIVE"),
	}
	title = dotg.EscapeMk(fmt.Sprintf("#测试 文件标题 %d", time.Now().UnixMilli()))
)

func TestFlvStream_Start(t *testing.T) {
	u := "https://cn-gddg-cu-01-04.bilivideo.com/live-bvc/660217/live_8739477_3713195_bluray.flv?expires=1681583774&pt=web&deadline=1681583774&len=0&oi=1885426465&platform=web&qn=10000&trid=10007de999a1b5e34d77be6b936cc001b229&uipk=100&uipv=100&nbs=1&uparams=cdn,deadline,len,oi,platform,qn,trid,uipk,uipv,nbs&cdn=cn-gotcha01&upsig=795bd27afbe3e15e0c454cbf04018c5c&sk=e3aba3a32d33f12ba204825894f385cb&p2p_type=0&sl=10&free_type=0&mid=0&sid=cn-gddg-cu-01-04&chash=0&sche=ban&score=18&pp=rtmp&source=one&trace=8a0&site=533ba65c94ba7a95dfb89645dca939f1&order=1"
	p := fmt.Sprintf("D:/Tmp/live/bili_6_%d.flv", time.Now().UnixMilli())
	s := NewFlvStream(title, u, plats.HeadersBili, p, 5*1024*1024, &tgHandler)

	err := s.Start()
	if err != nil {
		t.Fatal(err)
	}

	// 等待下载阶段的错误
	err = <-s.GetChErr()
	if err != nil {
		t.Fatal(err)
	}

	restart := <-s.GetChRestart()
	t.Logf("需要重新开始：%v", restart)
	if restart {
		TestFlvStream_Start(t)
	}
}
