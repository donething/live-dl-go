package stream

import (
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/live-dl-go/sites/douyin"
	"github.com/donething/live-dl-go/sites/entity"
	"github.com/donething/live-dl-go/sites/zuji"
	entity2 "github.com/donething/live-dl-go/stream/entity"
	"github.com/donething/live-dl-go/stream/entity/capture_status"
	"github.com/donething/utils-go/dotg"
	"os"
	"testing"
)

var (
	capturing = capture_status.New[entity2.IStream]()

	tgHandler = hanlders.TGHandler{
		TG:     dotg.NewTGBot(os.Getenv("MY_TG_TOKEN")),
		ChatID: os.Getenv("MY_TG_CHAT_LIVE"),
	}

	localHandle = hanlders.LocalHanlder{}
)

func TestStartAnchorFlv(t *testing.T) {
	anchor := entity.Anchor{
		ID:   "249406961231",
		Plat: douyin.Plat,
	}

	err := StartAnchor(capturing, nil, anchor, "D:/Tmp/live/douyin_249406961231.flv",
		10*1024*1024, &tgHandler)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartAnchorM3u8(t *testing.T) {
	anchor := entity.Anchor{
		ID:   "15722883",
		Plat: zuji.Plat,
	}

	err := StartAnchor(capturing, nil, anchor, "D:/Tmp/live/zuji_15722883.ts",
		5*1024*1024, &tgHandler)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartAnchorZuji(t *testing.T) {
	anchor := entity.Anchor{
		ID:   "20221998",
		Plat: zuji.Plat,
	}

	err := StartAnchor(capturing, nil, anchor, "D:/Tmp/live/zuji_20221998.ts",
		10*1024*1024, &localHandle)
	if err != nil {
		t.Fatal(err)
	}
}
