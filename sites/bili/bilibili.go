package bili

import (
	"encoding/json"
	"fmt"
	"live-dl-go/comm"
	"live-dl-go/sites/plats"
)

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

const tags = "[Bili]"

var headers = map[string]string{
	"user-agent": comm.UAWin,
}

// GetAnchorInfo 获取哔哩哔哩直播流的地址
//
// 参数 uid 为用户 ID，而不是房间号
func GetAnchorInfo(uid string) (*plats.AnchorInfo, error) {
	// 获取房间信息
	roomStatus, err := getRoomStatus(uid)
	if err != nil {
		return nil, err
	}
	playUrl, err := getPlayUrl(roomStatus.RoomID)
	if err != nil {
		return nil, err
	}

	anchor := plats.Anchor{
		ID:   uid,
		Plat: plats.PlatBili,
	}

	info := plats.AnchorInfo{
		Anchor:    &anchor,
		Avatar:    roomStatus.Face,
		Name:      roomStatus.Uname,
		WebUrl:    fmt.Sprintf("https://live.bilibili.com/%d", roomStatus.RoomID),
		Title:     roomStatus.Title,
		IsLive:    roomStatus.LiveStatus == 1,
		IsCycle:   roomStatus.LiveStatus == 2,
		StreamUrl: playUrl,
	}

	return &info, nil
}

// 获取哔哩哔哩直播间的信息
//
// 参数 uid 为用户 ID，而不是房间号
func getRoomStatus(uid string) (*RoomStatus, error) {
	url := "https://api.live.bilibili.com/room/v1/Room/get_status_info_by_uids?uids[]=" + uid
	bs, err := comm.Client.GetBytes(url, headers)
	if err != nil {
		return nil, fmt.Errorf("%s 获取主播信息出错：%w", tags, err)
	}

	var resp RespRoomStatus
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}

	// 因为直播地址对象的键是一个变量（为用户的 UID），不能直接定义
	// 需要先保证提取到该键对应的对象，再转为 RoomStatus 对象
	// 不通过断言转换是因为，RoomStatus 中只定义了部分属性，无法成功断言
	if resp.Data[uid] == nil {
		return nil, fmt.Errorf("无法定位到 UID 对应的直播间信息")
	}

	// 通过序列化、再反序列化，得到 RoomStatus
	tmp, err := json.Marshal(resp.Data[uid])
	if err != nil {
		return nil, err
	}
	var roomStatus RoomStatus
	err = json.Unmarshal(tmp, &roomStatus)
	if err != nil {
		return nil, err
	}

	return &roomStatus, nil
}

// 获取直播流地址
func getPlayUrl(roomid int) (string, error) {
	// 获取直播流地址
	url := fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/playUrl?platform=web&"+
		"qn=10000&cid=%d", roomid)
	bs, err := comm.Client.GetBytes(url, headers)
	if err != nil {
		return "", fmt.Errorf("%s 获取直播地址出错：%w", tags, err)
	}

	var resp RespPlayURL
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return "", fmt.Errorf("%s 解析直播地址出错：%w", tags, err)
	}

	// 解析得到地址
	if len(resp.Data.Durl) != 0 {
		return resp.Data.Durl[0].URL, nil
	}

	return "", fmt.Errorf("%s 直播流的地址数据为空", tags)
}
