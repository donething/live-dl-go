package stream

import (
	"testing"
)

func TestStream_StartM3u8(t *testing.T) {
	// localHandler := hanlders.LocalHanlder{}
	// title := "足迹"

	u := "xxx"
	p := "D:/Tmp/live/abc.ts"
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
