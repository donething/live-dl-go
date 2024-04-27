package stream

import (
	"fmt"
	"github.com/donething/live-dl-go/anchors/base"
	"github.com/donething/live-dl-go/anchors/bili"
	"github.com/donething/live-dl-go/anchors/douyin"
	"github.com/donething/live-dl-go/anchors/zuji"
	"github.com/donething/live-dl-go/hanlders"
	entity2 "github.com/donething/live-dl-go/stream/basestream"
	"github.com/donething/live-dl-go/stream/capture_status"
	"github.com/donething/utils-go/dotg"
	"os"
	"testing"
)

const workdir = "D:/Temp/VpsGo/anchors"

var (
	capturing = capture_status.New[entity2.IStream]()

	tgHandler = hanlders.TGHandler{
		TG:     dotg.NewTGBot(os.Getenv("MY_TG_TOKEN")),
		ChatID: os.Getenv("MY_TG_CHAT_LIVE"),
	}

	localHandle = hanlders.LocalHanlder{}

	afterHandle = func(task *hanlders.TaskInfo, err error) {
		if err != nil {
			fmt.Printf("处理视频出错(%s)：%s\n", task.Title, err)
			return
		}

		fmt.Printf("已处理视频(%s)到路径'%s'\n", task.Title, task.Path)
	}

	taskWithTG = hanlders.TaskInfo{
		Handler:           &tgHandler,
		FileSizeThreshold: 10 * 1024 * 1024,
		AfterHandle:       afterHandle,
	}

	taskWithLocal = hanlders.TaskInfo{
		Handler:           &localHandle,
		FileSizeThreshold: 10 * 1024 * 1024,
		AfterHandle:       afterHandle,
	}
)

func TestStartAnchorFlv(t *testing.T) {
	anchor := base.Anchor{
		UID:  "249406961231",
		Plat: douyin.Platform,
	}

	err := StartAnchor(capturing, anchor, workdir, taskWithTG)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartAnchorBili(t *testing.T) {
	anchor := base.Anchor{
		UID:  "8739477",
		Plat: bili.Platform,
	}

	err := StartAnchor(capturing, anchor, workdir, taskWithLocal)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartAnchorM3u8(t *testing.T) {
	anchor := base.Anchor{
		UID:  "15722883",
		Plat: zuji.Platform,
	}

	err := StartAnchor(capturing, anchor, workdir, taskWithTG)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartAnchorZuji(t *testing.T) {
	anchor := base.Anchor{
		UID:  "20221998",
		Plat: zuji.Platform,
	}

	err := StartAnchor(capturing, anchor, workdir, taskWithLocal)
	if err != nil {
		t.Fatal(err)
	}
}
