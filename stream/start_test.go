package stream

import (
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/live-dl-go/sites/bili"
	"github.com/donething/live-dl-go/sites/douyin"
	"github.com/donething/live-dl-go/sites/entity"
	"github.com/donething/live-dl-go/sites/zuji"
	"github.com/donething/utils-go/dotg"
	"os"
	"sync"
	"testing"
)

var (
	capturing = &sync.Map{}

	tgHandler = hanlders.TGHandler{
		TG:        dotg.NewTGBot(os.Getenv("MY_TG_TOKEN")),
		LocalPort: 0,
		ChatID:    os.Getenv("MY_TG_CHAT_LIVE"),
	}
)

// 哔哩哔哩的用户 ID
func TestStartFlv(t *testing.T) {
	anchor := entity.Anchor{
		ID:   "8739477",
		Plat: bili.Plat,
	}

	err := StartFlvAnchor(capturing, anchor, "D:/Tmp/live/bili_8739477.flv",
		20*1024*1024, &tgHandler)
	if err != nil {
		t.Fatal(err)
	}
}

// 抖音的直播间号
func TestStartFlv2(t *testing.T) {
	anchor := entity.Anchor{
		ID:   "249406961231",
		Plat: douyin.Plat,
	}

	err := StartFlvAnchor(capturing, anchor, "D:/Tmp/live/douyin_249406961231.flv",
		10*1024*1024, &tgHandler)
	if err != nil {
		t.Fatal(err)
	}
}

// 足迹的用户ID
func TestStartM3u8(t *testing.T) {
	anchor := entity.Anchor{
		ID:   "15722883",
		Plat: zuji.Plat,
	}

	err := StartM3u8Anchor(capturing, anchor, "D:/Tmp/live/zuji_61667788.flv",
		30*1024*1024, &tgHandler)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartAnchor(t *testing.T) {
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
