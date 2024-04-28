package bili

// RespRoomStatus 获取直播间的信息
type RespRoomStatus struct {
	// 0：成功；-111：csrf校验失败
	Code    int    `json:"code"`
	Message string `json:"message"`
	// data[UID] 为 RoomStatus 的实例
	Data map[string]interface{} `json:"data"`
}
type RoomStatus struct {
	Uname  string `json:"uname"`
	Title  string `json:"title"`
	RoomID int    `json:"room_id"`
	// 0：未开播；1：正在直播；2：轮播中
	LiveStatus int `json:"live_status"`
	// 头像
	Face string `json:"face"`
}

// RespPlayURL 获取直播流
type RespPlayURL struct {
	// 0：成功；-400：参数错误；19002003：房间信息不存在
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		// 直播流
		Durl []struct {
			// flv 或 m3u8 格式地址。带有 Unicode 转义字符，如"\u1234"
			URL string `json:"url"`
		} `json:"durl"`
	} `json:"data"`
}
