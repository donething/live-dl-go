package bili

import (
	"encoding/json"
	"fmt"
	"github.com/donething/live-dl-go/anchors/base"
	"github.com/donething/live-dl-go/request"
)

// AnchorBili 哔哩哔哩主播
type AnchorBili struct {
	// 主播的 UID 为用户 ID，非房间号
	*base.Anchor
}

const (
	Platform = "bili"
	Name     = "哔哩哔哩"
)

// GetAnchorInfo 获取哔哩哔哩直播流的地址
func (a *AnchorBili) GetAnchorInfo() (*base.AnchorInfo, error) {
	// 获取房间信息
	roomStatus, err := getRoomStatus(a.UID)
	if err != nil {
		return base.GenAnchorInfoWhenErr(a.Anchor, fmt.Sprintf("https://space.bilibili.com/%s", a.UID)), err
	}
	playUrl, err := getPlayUrl(roomStatus.RoomID)
	if err != nil {
		return base.GenAnchorInfoWhenErr(a.Anchor, fmt.Sprintf("https://space.bilibili.com/%s", a.UID)), err
	}

	info := base.AnchorInfo{
		Anchor:    a.Anchor,
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

// GetPlatName 获取平台名
func (a *AnchorBili) GetPlatName() string {
	return Name
}

// GetStreamHeaders 请求直播流时的请求头
func (a *AnchorBili) GetStreamHeaders() map[string]string {
	return map[string]string{
		// referer 必不可少
		"referer":    "https://live.bilibili.com/",
		"user-agent": request.UAWin,
	}
}

// 获取哔哩哔哩直播间的信息
func getRoomStatus(uid string) (*RoomStatus, error) {
	headers := map[string]string{
		"user-agent": request.UAWin,
	}
	url := "https://api.live.bilibili.com/room/v1/Room/get_status_info_by_uids?uids[]=" + uid
	bs, err := request.Client.GetBytes(url, headers)
	if err != nil {
		return nil, fmt.Errorf("获取主播信息出错：%w", err)
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
	headers := map[string]string{
		"user-agent": request.UAWin,
	}
	// 获取直播流地址
	url := fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/playUrl?platform=web&"+
		"qn=10000&cid=%d", roomid)
	bs, err := request.Client.GetBytes(url, headers)
	if err != nil {
		return "", fmt.Errorf("获取直播地址出错：%w", err)
	}

	var resp RespPlayURL
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return "", fmt.Errorf("解析直播地址出错：%w", err)
	}

	// 解析得到地址
	if len(resp.Data.Durl) != 0 {
		return resp.Data.Durl[0].URL, nil
	}

	return "", fmt.Errorf("直播流的地址数据为空")
}
