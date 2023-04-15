// Package douyin 获取抖音直播间的状态
// [douyin_dynamic_push](https://github.com/nfe-w/douyin_dynamic_push/blob/master/query_douyin.py)
// [Python爬取抖音用户相关数据(目前最方便的方法）](http://www.dagoogle.cn/n/1307.html)
package douyin

import (
	"dl-live-go/comm"
	"dl-live-go/sites/plats"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
)

var (
	// 目前除了要浏览器代理，还需要提供 cookie，否则获取到的是滑动验证页面
	headers = map[string]string{
		"user-agent": comm.UAWin,
		"referer":    "https://live.douyin.com/",
		"cookie":     `__ac_nonce=0643858c700f34a8c2c8a; __ac_signature=_02B4Z6wo00f01P5n5.QAAIDBnW0nnZXIxDD-R-NAAFun3d; ttwid=1%7Cel-5Doi29H-jlz0Xx2hq2NIGFu52ojyovikjEvWE7nI%7C1681414343%7Ca8d601d269cb3a2d959fbc84f1efbfc12f994878c9d218f8840730fc4984a95a; __live_version__=%221.1.0.8512%22; device_web_cpu_core=6; device_web_memory_size=8; live_can_add_dy_2_desktop=%220%22; csrf_session_id=80cbc87883f90a294efa5df9d1d1071f; webcast_local_quality=sd; xgplayer_user_id=420141080552; odin_tt=70b997c3ad3c3d2739734731de4414c1557629c28e7805699d211b657bd8b9af64507faaffe7ed73bb3c17921b343a72f9f98f98d434ff5ff667d32bbdc9b72addd26f309ab3eade2f4a825574b0cf71; s_v_web_id=verify_lgfiovrb_Qfn1H8Or_tVGT_40CI_9ICx_OfJcZyfxNrXL; passport_csrf_token=200310d76f2b1f67874730055d6eeb05; passport_csrf_token_default=200310d76f2b1f67874730055d6eeb05; msToken=Wx0wVzU4c_4Xqtps0BxewKzNqdUdlh4KK7cr207MEYSljxBEQ-ZQNLz7ewkdVDY9bpBe5IWCo12ronkhBxxXU6gFNu05g9ZoTOy8wmQzV8GEE1s8PHZjzUfYSfDiwA==; ttcid=96084af587fd43a09f67e0a2817cb0c118; tt_scid=K4G-bpw79V9jdx1FA5ciYmyxnEA.9mLKUMusCHBUE.aQfkK010CR.Ag78hx2.L1P7f86; msToken=h6I_i5D3SH7DkA9jDV2X3UX3qYhAu8FLqM91h0yIGdKyWLF3KTJEZc6-WOk3P_KeATd_TaaOJEUCghJ4ugXs9CoCif14r5f0l4gXCXjcwkHAan7F2skbo1EdNK7q1w==`,
	}
)

// GetAnchorInfo 获取抖音主播直播间的信息
//
// roomid 直播间号
func GetAnchorInfo(roomid string) (*plats.AnchorInfo, error) {
	// 提取直播间的直播信息
	u := fmt.Sprintf("https://live.douyin.com/%s", roomid)
	roomStatus, err := parseRenderData[RoomStatus](u)
	if err != nil {
		return nil, fmt.Errorf("获取直播间信息出错(%s)：%w", roomid, err)
	}

	// 是否开播，关系到页面中是否存在数据
	if roomStatus.App.InitialState.RoomStore.RoomInfo.Anchor.Nickname == "" {
		return nil, fmt.Errorf("不存在直播间(%s)", roomid)
	}

	anchor := plats.Anchor{
		ID:   roomid,
		Plat: plats.PlatDouyin,
	}

	roomInfo := roomStatus.App.InitialState.RoomStore.RoomInfo
	anchorInfo := plats.AnchorInfo{
		Anchor:    &anchor,
		Avatar:    roomInfo.Anchor.AvatarThumb.URLList[0],
		Name:      roomInfo.Anchor.Nickname,
		WebUrl:    fmt.Sprintf("https://live.douyin.com/%s", roomInfo.WebRid),
		Title:     roomInfo.Room.Title,
		IsLive:    roomStatus.App.InitialState.RoomStore.RoomInfo.Room.Status == 2,
		StreamUrl: roomInfo.Room.StreamURL.FlvPullURL.FULLHD1,
	}

	return &anchorInfo, nil
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
