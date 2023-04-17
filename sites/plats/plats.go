package plats

import (
	"fmt"
	"github.com/donething/live-dl-go/sites/bili"
	"github.com/donething/live-dl-go/sites/douyin"
	"github.com/donething/live-dl-go/sites/entity"
	"github.com/donething/live-dl-go/sites/zuji"
)

// Plats 已适配的所有平台
var Plats = map[string]string{
	bili.Plat:   "",
	douyin.Plat: "",
	zuji.Plat:   "",
}

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

// ExistPlat 检查给定的平台是否存在
func ExistPlat(plat string) bool {
	_, exists := Plats[plat]
	return exists
}
