package flv

import (
	"fmt"
	"github.com/donething/live-dl-go/anchors/base"
	"github.com/donething/live-dl-go/anchors/bili"
	"github.com/donething/live-dl-go/anchors/platform"
	"github.com/donething/live-dl-go/hanlders"
	"path/filepath"
	"testing"
)

const workdir = "D:/Temp/VpsGo/anchors"

var (
	localHandle = hanlders.LocalHanlder{}

	afterHandle = func(task *hanlders.TaskInfo, err error) {
		if err != nil {
			fmt.Printf("处理视频出错(%s)：%s\n", task.Title, err)
			return
		}

		fmt.Printf("已处理视频(%s)到路径'%s'\n", task.Title, task.Path)
	}

	taskWithLocal = hanlders.TaskInfo{
		Path:              filepath.Join(workdir, "测试 文件标题.flv"),
		Title:             "测试 文件标题.flv",
		Handler:           &localHandle,
		FileSizeThreshold: 10 * 1024 * 1024,
		AfterHandle:       afterHandle,
	}
)

func TestStream_Capture(t *testing.T) {
	anchor := base.Anchor{
		UID:  "2011822166",
		Plat: bili.Platform,
	}

	anchorSite, err := platform.GenAnchor(&anchor)
	if err != nil {
		t.Fatal(err)
	}

	s := NewStream(&taskWithLocal, anchorSite)
	err = s.Capture()
	if err != nil {
		t.Fatal(err)
	}
}
