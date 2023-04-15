package stream

import (
	"fmt"
	"testing"
	"time"
)

func TestStream_StartM3u8(t *testing.T) {
	// localHandler := hanlders.LocalHanlder{}
	// title := "足迹"

	u := "http://bjlive.szsbtech.com/record/ApK9fkvPjVrCP6x.m3u8?auth_key=1681561481-0-0-844400a36b4decd30b35d3fb95d0ff33"
	p := fmt.Sprintf("D:/Tmp/live/zuji_%d.ts", time.Now().UnixMilli())
	s := NewM3u8Stream(title, u, nil, p, 5*1024*1024, &tgHandler)

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
