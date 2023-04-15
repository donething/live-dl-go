package zuji

import (
	"testing"
)

func TestGetAnchorInfo(t *testing.T) {
	// 放在播的
	info, err := GetAnchorInfo("20225898")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", info)

	// 放不在播的
	info, err = GetAnchorInfo("15050303")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", info)
}
