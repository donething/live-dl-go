package zuji

// RespInterface 1. 用来获取直播间号
//
// 本来这个响应内容基本和 RespWatchStart 一直。但是缺少最关键的 `play_url` (有属性，但为空"")
//
// 不得不根据这个获取 vid，再获取 RespWatchStart 得到 play_url
type RespInterface struct {
	// "ok" 表示成功，其它字符串表示出错
	Retval string `json:"retval"`
	Reterr string `json:"reterr"`

	Retinfo struct {
		// 1 表示在播，0 表示未开播或为付费房
		Roomstatus int `json:"roomstatus"`
		// 临时直播间号。用于获取直播流地址
		Vid string `json:"vid"`
		// 直播间标题
		Title string `json:"title"`
		// 用户显示的 ID。如短号、靓号
		Name string `json:"name"`
		// 用户真实 ID
		Uid int `json:"uid"`
		// 用户昵称
		Nickname string `json:"nickname"`
		// 用户头像
		Logourl string `json:"logourl"`
		// 直播间的权限。0 免费；7 收费
		Permission int `json:"permission"`
	} `json:"retinfo"`
}

// RespWatchStart 2. 获取直播流地址
//
// GET https://m.szsbtech.com/appgw/v2/watchstart
type RespWatchStart struct {
	// "ok" 表示成功，其它字符串表示出错
	Retval string `json:"retval"`
	// 可读的描述。可能出现出错，但此值为空的情况。建议用 Retval 显示出错
	Reterr string `json:"reterr"`

	Retinfo struct {
		// 直播流地址。Android UA 返回 `m3u8` 流，而 PC 端返回 "rtmp://"流
		// 可以模拟 Android UA 获取 `m3u8`直播流
		PlayURL string `json:"play_url"`
	} `json:"retinfo"`
}
