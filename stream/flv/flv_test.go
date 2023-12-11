package flv

import (
	"fmt"
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/live-dl-go/sites/bili"
	"github.com/donething/live-dl-go/sites/entity"
	"github.com/donething/live-dl-go/sites/plats"
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
	title = dotg.EscapeMk(fmt.Sprintf("#测试 文件标题 %d", time.Now().Unix()))
)

func TestStream_Capture(t *testing.T) {
	anchor := entity.Anchor{
		UID:  "2011822166",
		Plat: bili.Plat,
	}

	anchorSite, err := plats.GenAnchor(&anchor)
	if err != nil {
		t.Fatal(err)
	}

	s := NewStream(title, "D:/Temp/VpsGo/captures/bili_616.flv", 5*1024*1024, &tgHandler, anchorSite)
	err = s.Capture()
	if err != nil {
		t.Fatal(err)
	}
}
