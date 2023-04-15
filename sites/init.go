package sites

import (
	"dl-live-go/sites/bili"
	"dl-live-go/sites/douyin"
	"dl-live-go/sites/plats"
	"dl-live-go/sites/zuji"
)

func init() {
	// 初始化各个平台的操作
	plats.Plats[plats.PlatBili] = &plats.PlatOp{GetAnchorInfo: bili.GetAnchorInfo}
	plats.Plats[plats.PlatDouyin] = &plats.PlatOp{GetAnchorInfo: douyin.GetAnchorInfo}
	plats.Plats[plats.PlatZuji] = &plats.PlatOp{GetAnchorInfo: zuji.GetAnchorInfo}
}
