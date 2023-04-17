package plats

import (
	"fmt"
	"github.com/donething/live-dl-go/sites/bili"
	"github.com/donething/live-dl-go/sites/douyin"
	"github.com/donething/live-dl-go/sites/entity"
	"github.com/donething/live-dl-go/sites/zuji"
)

// GenAnchor 自动生成与平台对应的 Anchor* 的实例
func GenAnchor(anchor *entity.Anchor) (entity.IAnchor, error) {
	switch anchor.Plat {
	case bili.Plat:
		return &bili.AnchorBili{Anchor: anchor}, nil
	case douyin.Plat:
		return &douyin.AnchorDouyin{Anchor: anchor}, nil
	case zuji.Plat:
		return &zuji.AnchorZuji{Anchor: anchor}, nil
	default:
		return nil, fmt.Errorf("未知的平台(%+v)", anchor)
	}
}
