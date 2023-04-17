// Package douyin 获取抖音直播间的状态
// [douyin_dynamic_push](https://github.com/nfe-w/douyin_dynamic_push/blob/master/query_douyin.py)
// [Python爬取抖音用户相关数据(目前最方便的方法）](http://www.dagoogle.cn/n/1307.html)
package douyin

import (
	"encoding/json"
	"fmt"
	"github.com/donething/live-dl-go/comm"
	"github.com/donething/live-dl-go/sites/entity"
	"net/url"
	"regexp"
)

// AnchorDouyin 抖音主播
type AnchorDouyin struct {
	// 主播的 ID 为直播间号
	*entity.Anchor
}

const (
	Plat = "douyin"
	name = "抖音"
)

var (
	// 目前除了要浏览器代理，还需要提供 cookie，否则获取到的是滑动验证页面
	headers = map[string]string{
		"user-agent": comm.UAWin,
		"referer":    "https://live.douyin.com/",
		"cookie":     `__ac_nonce=0643d7054003a052ad132; __ac_signature=_02B4Z6wo00f01LAm0.wAAIDB0ywTl16l12SwBtdAAEg5ad; __ac_referer=__ac_blank`,
	}
)

// GetAnchorInfo 获取抖音主播直播间的信息
//
// roomid 直播间号
func (a *AnchorDouyin) GetAnchorInfo() (*entity.AnchorInfo, error) {
	// 提取直播间的直播信息
	u := fmt.Sprintf("https://live.douyin.com/%s", a.ID)
	roomStatus, err := parseRenderData[RoomStatus](u)
	if err != nil {
		return entity.GenAnchorInfoWhenErr(a.Anchor, fmt.Sprintf("https://live.douyin.com/%s", a.ID)),
			fmt.Errorf("获取直播间出错：%w", err)
	}

	// 是否开播，关系到页面中是否存在数据
	if roomStatus.App.InitialState.RoomStore.RoomInfo.Anchor.Nickname == "" {
		return entity.GenAnchorInfoWhenErr(a.Anchor, fmt.Sprintf("https://live.douyin.com/%s", a.ID)),
			fmt.Errorf("不存在的直播间")
	}

	roomInfo := roomStatus.App.InitialState.RoomStore.RoomInfo
	anchorInfo := entity.AnchorInfo{
		Anchor:    a.Anchor,
		Avatar:    roomInfo.Anchor.AvatarThumb.URLList[0],
		Name:      roomInfo.Anchor.Nickname,
		WebUrl:    fmt.Sprintf("https://live.douyin.com/%s", roomInfo.WebRid),
		Title:     roomInfo.Room.Title,
		IsLive:    roomStatus.App.InitialState.RoomStore.RoomInfo.Room.Status == 2,
		StreamUrl: roomInfo.Room.StreamURL.FlvPullURL.FULLHD1,
	}

	return &anchorInfo, nil
}

// GetPlatName 获取平台名
func (a *AnchorDouyin) GetPlatName() string {
	return name
}

// GetStreamHeaders 请求直播流时的请求头
func (a *AnchorDouyin) GetStreamHeaders() map[string]string {
	return map[string]string{
		"user-agent": comm.UAWin,
	}
}

// 提取网页中的 RENDER_DATA
func parseRenderData[T any](dyUrl string) (*T, error) {
	// 获取抖音网页文本
	text, err := comm.Client.GetText(dyUrl, headers)
	if err != nil {
		return nil, fmt.Errorf("获取网页内容出错(%s)：%w", dyUrl, err)
	}

	// 页面会携带一段ID为"RENDER_DATA"的脚本，里面带有用户数据信息
	// 可以在页面控制台中执行`copy(decodeURIComponent(document.querySelector("#RENDER_DATA").text))`获取
	reg := regexp.MustCompile(`(?m)id="RENDER_DATA".+?>(.+?)<`)
	matches := reg.FindStringSubmatch(text)
	if len(matches) < 2 {
		return nil, fmt.Errorf("没有匹配到'RENDER_DATA'数据")
	}

	// 反转义非法字符
	dataText, err := url.QueryUnescape(matches[1])
	if err != nil {
		return nil, fmt.Errorf("反转义非法字符出错：%w", err)
	}

	// 解析数据
	var obj = new(T)
	err = json.Unmarshal([]byte(dataText), obj)

	if err != nil {
		return nil, fmt.Errorf("解析数据出错：%w", err)
	}

	return obj, nil
}
