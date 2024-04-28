package platform

import (
	"fmt"
	"github.com/donething/live-dl-go/anchors/baseanchor"
	"github.com/donething/live-dl-go/anchors/sites/bili"
	"github.com/donething/live-dl-go/anchors/sites/douyin"
	"github.com/donething/live-dl-go/anchors/sites/zuji"
)

// Platforms 已适配的所有平台
var Platforms = map[string]string{
	bili.Platform:   "",
	douyin.Platform: "",
	zuji.Platform:   "",
}

// GenAnchor 自动生成与平台对应的 Anchor* 的实例
func GenAnchor(anchor *baseanchor.Anchor) (baseanchor.IAnchor, error) {
	switch anchor.Plat {
	case bili.Platform:
		return &bili.AnchorBili{Anchor: anchor}, nil

	case douyin.Platform:
		return &douyin.AnchorDouyin{Anchor: anchor}, nil

	case zuji.Platform:
		return &zuji.AnchorZuji{Anchor: anchor}, nil

	default:
		return nil, fmt.Errorf("未知的平台(%+v)", anchor)
	}
}

// ExistPlat 检查给定的平台是否存在
func ExistPlat(plat string) bool {
	_, exists := Platforms[plat]
	return exists
}
