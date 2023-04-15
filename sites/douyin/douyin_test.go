package douyin

import (
	"testing"
)

func TestGetLiveStream(t *testing.T) {
	liveInfo, err := GetAnchorInfo("165251594775")
	if err != nil {
		t.Log(err)
		return
	}

	t.Logf("%+v\n", liveInfo)
}
