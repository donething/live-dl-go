package sites

import (
	"github.com/donething/live-dl-go/sites/bili"
	"github.com/donething/live-dl-go/sites/douyin"
	"github.com/donething/live-dl-go/sites/plats"
	"github.com/donething/live-dl-go/sites/zuji"
)

func init() {
	// 初始化各个平台的操作
	plats.Plats[plats.PlatBili] = &plats.PlatOp{GetAnchorInfo: bili.GetAnchorInfo}
	plats.Plats[plats.PlatDouyin] = &plats.PlatOp{GetAnchorInfo: douyin.GetAnchorInfo}
	plats.Plats[plats.PlatZuji] = &plats.PlatOp{GetAnchorInfo: zuji.GetAnchorInfo}
}
