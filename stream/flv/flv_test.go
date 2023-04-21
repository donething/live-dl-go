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

func TestFlvStream_Start(t *testing.T) {
	headers := map[string]string{
		// referer 必不可少
		"referer":    "https://live.bilibili.com/",
		"user-agent": comm.UAWin,
	}
	u := "https://cn-gddg-cu-01-04.bilivideo.com/live-bvc/339767/live_8739477_3713195_bluray.flv?expires=1681755864&pt=web&deadline=1681755864&len=0&oi=1885426465&platform=web&qn=10000&trid=1000c4a0a4b4d2d648ae86e9eff7342030c7&uipk=100&uipv=100&nbs=1&uparams=cdn,deadline,len,oi,platform,qn,trid,uipk,uipv,nbs&cdn=cn-gotcha01&upsig=f0ff7706bd14896b1115ec4c17fe4049&sk=d3424ef93a341b3317ba4806163c6ba8&p2p_type=0&sl=10&free_type=0&mid=0&sid=cn-gddg-cu-01-04&chash=0&sche=ban&score=18&pp=rtmp&source=one&trace=8a0&site=6cd4239d0deb0e46a732c0a7fc976b75&order=1"
	p := fmt.Sprintf("D:/Tmp/live/bili_6_%d.flv", time.Now().UnixMilli())
	s := NewStream(title, u, headers, p, 5*1024*1024, &tgHandler)

	err := s.Start()
	if err != nil {
		t.Fatal(err)
	}

	// 等待下载阶段的错误
	err = <-s.GetStream().ChErr
	if err != nil {
		t.Fatal(err)
	}

	restart := <-s.GetStream().ChRestart
	t.Logf("需要重新开始：%v", restart)
	if restart {
		TestFlvStream_Start(t)
	}
}

func TestStream_Download(t *testing.T) {
	headers := map[string]string{
		// referer 必不可少
		"referer":    "https://live.bilibili.com/",
		"user-agent": comm.UAWin,
	}
	u := "https://aaa.com/downloads/uploads/01.mp4"
	p := fmt.Sprintf("D:/Tmp/live/%d_video.mp4", time.Now().UnixMilli())

	err := Download(title, u, headers, p, &tgHandler)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Second)
}
