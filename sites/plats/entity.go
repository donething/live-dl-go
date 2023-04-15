package plats

// Anchor 需要获取直播流的主播的信息
type Anchor struct {
	// 直播流
	ID string `json:"id"`
	// 直播平台。必须为`comm.go`中的常量 Plat** 项
	Plat string `json:"plat"`
}

// AnchorInfo 直播房间的信息
type AnchorInfo struct {
	*Anchor
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Desp   string `json:"desp"`
	// Web 端观看地址
	WebUrl string `json:"webUrl"`

	// 直播间的标题
	Title string `json:"title"`
	// 是否开播
	IsLive bool `json:"isLive"`
	// 是否轮播
	IsCycle bool `json:"isCycle"`
	// 直播流的地址
	StreamUrl string `json:"streamUrl"`
}

// PlatOp 不同平台对应的操作
type PlatOp struct {
	GetAnchorInfo func(uid string) (*AnchorInfo, error)
}
