package sites

import (
	"live-dl-go/sites/bili"
	"live-dl-go/sites/douyin"
	"live-dl-go/sites/plats"
	"live-dl-go/sites/zuji"
)

func init() {
	// 初始化各个平台的操作
	plats.Plats[plats.PlatBili] = &plats.PlatOp{GetAnchorInfo: bili.GetAnchorInfo}
	plats.Plats[plats.PlatDouyin] = &plats.PlatOp{GetAnchorInfo: douyin.GetAnchorInfo}
	plats.Plats[plats.PlatZuji] = &plats.PlatOp{GetAnchorInfo: zuji.GetAnchorInfo}
}
