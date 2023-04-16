package plats

import "github.com/donething/live-dl-go/comm"

// 平台
const (
	PlatBili   = "bili"
	PlatDouyin = "douyin"
	PlatZuji   = "zuji"
)

// Sites 平台对应的网站名
var Sites = map[string]string{
	PlatBili:   "哔哩哔哩",
	PlatDouyin: "抖音",
	PlatZuji:   "足迹",
}

// Plats 主播的平台
//
// 避免循环导入包，此处不能直接初始化
//
// 务必在 `live/web.go` 的 `init()` 中对各个平台完成赋值
var Plats = map[string]*PlatOp{
	PlatBili:   {},
	PlatDouyin: {},
	PlatZuji:   {},
}

// HeadersBili 哔哩哔哩直播的请求头
var HeadersBili = map[string]string{
	// referer 必不可少
	"referer":    "https://live.bilibili.com/",
	"user-agent": comm.UAWin,
}

var HeadersComm = map[string]string{
	"user-agent": comm.UAWin,
}
