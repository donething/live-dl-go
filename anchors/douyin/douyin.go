// Package douyin 获取抖音直播间的状态
// [douyin_dynamic_push](https://github.com/nfe-w/douyin_dynamic_push/blob/master/query_douyin.py)
// [Python爬取抖音用户相关数据(目前最方便的方法）](http://www.dagoogle.cn/n/1307.html)
package douyin

import (
	"encoding/json"
	"fmt"
	"github.com/donething/live-dl-go/anchors/base"
	"github.com/donething/live-dl-go/request"
	"net/url"
	"regexp"
)

// AnchorDouyin 抖音主播
type AnchorDouyin struct {
	// 主播的 UID 为直播间号
	*base.Anchor
}

const (
	Platform = "douyin"
	Name     = "抖音"
)

var (
	// 目前除了要浏览器代理，还需要提供 cookie，否则获取到的是滑动验证页面
	headers = map[string]string{
		"user-agent": request.UAWin,
		"referer":    "https://live.douyin.com/",
		"cookie": "device_web_cpu_core=6; device_web_memory_size=8; xgplayer_user_id=954101977251; " +
			"csrf_session_id=f7a065a09fe1ab43fb54d0912daca6bd; webcast_leading_last_show_time=1693924374823; " +
			"webcast_leading_total_show_times=1; webcast_local_quality=origin; store-region=cn-gd; " +
			"store-region-src=uid; LOGIN_STATUS=1; d_ticket=83146d33dfd9d7ce477f7b5302529d01b3913; " +
			"my_rd=2; live_use_vvc=%22false%22; passport_csrf_token=27e62cc168b78210e0d1855b4760d2de; " +
			"passport_csrf_token_default=27e62cc168b78210e0d1855b4760d2de; " +
			"bd_ticket_guard_client_data=eyJiZC10aWNrZXQtZ3VhcmQtdmVyc2lvbiI6MiwiYmQtdGlja2V0LWd1YXJkLWl" +
			"0ZXJhdGlvbi12ZXJzaW9uIjoxLCJiZC10aWNrZXQtZ3VhcmQtcmVlLXB1YmxpYy1rZXkiOiJCR2dBVFBybmE3NmxpYi9Y" +
			"MnpUL3N1SFJ4NTh3dFdad2FQY2RmdTRWZU4yck1TT1F6Z3lJQlpKK2EwQkljNUFuQXRuWFcrNnRhVDRDKzhJVzAvWWRscm" +
			"s9IiwiYmQtdGlja2V0LWd1YXJkLXdlYi12ZXJzaW9uIjoxfQ%3D%3D; passport_fe_beating_status=false; " +
			"bd_ticket_guard_client_web_domain=2; n_mh=SQMIVkn0ZI3dBLrve6IKic_UJ7rK0bsMaYj6timX0J8; " +
			"passport_auth_status=c5b6c1a3a796abc0632b8e729e468b8a%2C; passport_auth_status_ss=c5b6c1a3a" +
			"796abc0632b8e729e468b8a%2C; _bd_ticket_crypt_doamin=2; _bd_ticket_crypt_cookie=ab8022a40e9100e2" +
			"d88e50f846c07056; __security_server_data_status=1; ttwid=1%7Cw1fq-7nNm5SobrfvxPsT-sExiUFR_HbhZzJ" +
			"GCKpTtX4%7C1701814122%7C066011984b25362a713bcba4794a5e2a7dbd36ea1e9a7297ca50522280466dbd; " +
			"FORCE_LOGIN=%7B%22videoConsumedRemainSeconds%22%3A180%7D; __live_version__=%221.1.1.6369%22;" +
			" webcast_local_quality=sd; download_guide=%223%2F20231218%2F0%22; pwa2=%220%7C0%7C3%7C0%22; " +
			"strategyABtestKey=%221702894453.983%22; sid_guard=d330259d2d4a279aef49f1d7b227364e%7C17028944" +
			"54%7C21600%7CMon%2C+18-Dec-2023+16%3A14%3A14+GMT; uid_tt=8db379233d4651891ae98679bbfd296b; " +
			"uid_tt_ss=8db379233d4651891ae98679bbfd296b; sid_tt=d330259d2d4a279aef49f1d7b227364e; " +
			"sessionid=d330259d2d4a279aef49f1d7b227364e; sessionid_ss=d330259d2d4a279aef49f1d7b227364e; " +
			"sid_ucp_v1=1.0.0-KDVjYzAxNjc0NTIwNDNkNDkyMTg0Njc0Yjk1NmExMzQzNjJiM2Y1OTQKCBD2toCsBhgNGgJob" +
			"CIgZDMzMDI1OWQyZDRhMjc5YWVmNDlmMWQ3YjIyNzM2NGU; ssid_ucp_v1=1.0.0-KDVjYzAxNjc0NTIwNDNkNDkyMTg0" +
			"Njc0Yjk1NmExMzQzNjJiM2Y1OTQKCBD2toCsBhgNGgJobCIgZDMzMDI1OWQyZDRhMjc5YWVmNDlmMWQ3YjIyNzM2NGU; " +
			"volume_info=%7B%22isUserMute%22%3Afalse%2C%22isMute%22%3Afalse%2C%22volume%22%3A1%7D; " +
			"odin_tt=cb883421ae36c36c26657024e478ec8fa8b3b65159e7bfce38308a8ced85466096dd955389cd89e02c" +
			"d5a26ee18e28db122affef47426960985d0f6f916bbf74ab0f5c0a2d6fb5c69f45a6fe7de24f6e; " +
			"xg_device_score=7.345736471131223; bd_ticket_guard_client_data=eyJiZC10aWNrZXQtZ3VhcmQtdmVyc" +
			"2lvbiI6MiwiYmQtdGlja2V0LWd1YXJkLWl0ZXJhdGlvbi12ZXJzaW9uIjoxLCJiZC10aWNrZXQtZ3VhcmQtcmVlLXB" +
			"1YmxpYy1rZXkiOiJCR2dBVFBybmE3NmxpYi9YMnpUL3N1SFJ4NTh3dFdad2FQY2RmdTRWZU4yck1TT1F6Z3lJQlpKK2EwQ" +
			"kljNUFuQXRuWFcrNnRhVDRDKzhJVzAvWWRscms9IiwiYmQtdGlja2V0LWd1YXJkLXdlYi12ZXJzaW9uIjoxfQ%3D%3D; " +
			"stream_player_status_params=%22%7B%5C%22is_auto_play%5C%22%3A0%2C%5C%22is_full_screen%5C%22%3A0%2" +
			"C%5C%22is_full_webscreen%5C%22%3A1%2C%5C%22is_mute%5C%22%3A0%2C%5C%22is_speed%5C%22%3A1%2C%" +
			"5C%22is_visible%5C%22%3A1%7D%22; stream_recommend_feed_params=%22%7B%5C%22cookie_enabled%5C%22" +
			"%3Atrue%2C%5C%22screen_width%5C%22%3A1920%2C%5C%22screen_height%5C%22%3A1080%2C%5C%22browser_on" +
			"line%5C%22%3Atrue%2C%5C%22cpu_core_num%5C%22%3A6%2C%5C%22device_memory%5C%22%3A8%2C%5C%22downli" +
			"nk%5C%22%3A10%2C%5C%22effective_type%5C%22%3A%5C%224g%5C%22%2C%5C%22round_trip_time%5C%22%3A150%7D%22;" +
			" __ac_nonce=065802c590010beaea26c; __ac_signature=_02B4Z6wo00f01-lGOcQAAIDC42DibhQ9r7PpZj1AAJ.bi8D" +
			"lUph1vl.J88-ogCLvrVAMGrPhbvWmD2FVu-VfadGBm8dXl1i0Xb1irw9MD6rkr85YPFb40s6THIuf4x3gJQmf5xcvmgSowWOO56; " +
			"webcast_leading_last_show_time=1702898780349; webcast_leading_total_show_times=19; " +
			"live_can_add_dy_2_desktop=%221%22; msToken=DTvtwPAsS2H_9Oqn-U1JE48En1cHYEfzQptIpYKGcHXapLzHT4tZweNv1" +
			"nkNQcH_5isKcloRLmLp66qCTLS6sqx4SdMqUNeB0Boc19SH_Gj2oJqnfaR2MD3-bkwq; tt_scid=utwOYf6O1GBmKER7JRLcZClV" +
			"djhLtH8ZbQupmAHoAjWRI966-b7xs1uhRz0tnyGR450a; home_can_add_dy_2_desktop=%221%22; IsDouyinActive=false; ",
	}
)

// GetAnchorInfo 获取抖音主播直播间的信息
func (a *AnchorDouyin) GetAnchorInfo() (*base.AnchorInfo, error) {
	// 提取直播间的直播信息
	// 此 API 请求头需要 Cookie，可以为未登录时的 Cookie
	u := fmt.Sprintf("https://live.douyin.com/webcast/room/web/enter/?aid=6383&device_platform=web&"+
		"enter_from=web_live&cookie_enabled=true&browser_language=zh-CN&browser_platform=Win32&"+
		"browser_name=Chrome&browser_version=109.0.0.0&web_rid=%s", a.UID)
	bs, err := request.Client.GetBytes(u, headers)
	if err != nil {
		return base.GenAnchorInfoWhenErr(a.Anchor, fmt.Sprintf("https://live.douyin.com/%s", a.UID)),
			fmt.Errorf("获取直播间出错：%w", err)
	}

	// 解析数据
	var obj RoomInfo
	err = json.Unmarshal(bs, &obj)
	if err != nil {
		return nil, fmt.Errorf("解析数据出错：%w", err)
	}

	if obj.StatusCode != 0 {
		return nil, fmt.Errorf("主播不存在")
	}

	roomInfo := obj.Data.Data[0]
	anchorInfo := base.AnchorInfo{
		Anchor:    a.Anchor,
		Avatar:    obj.Data.User.AvatarThumb.URLList[0],
		Name:      obj.Data.User.Nickname,
		WebUrl:    fmt.Sprintf("https://www.douyin.com/user/%s", roomInfo.Owner.SecUID),
		Title:     roomInfo.Title,
		IsLive:    obj.Data.Data[0].StreamURL.HlsPullURL != "",
		StreamUrl: roomInfo.StreamURL.HlsPullURL,
	}

	return &anchorInfo, nil
}

// GetAnchorInfoParseWeb 获取抖音主播直播间的信息
//
// roomid 直播间号
func (a *AnchorDouyin) GetAnchorInfoParseWeb() (*base.AnchorInfo, error) {
	// 提取直播间的直播信息
	u := fmt.Sprintf("https://live.douyin.com/%s", a.UID)
	roomStatus, err := parseRenderData[RoomStatus](u)
	if err != nil {
		return base.GenAnchorInfoWhenErr(a.Anchor, fmt.Sprintf("https://live.douyin.com/%s", a.UID)),
			fmt.Errorf("获取直播间出错：%w", err)
	}

	// 是否开播，关系到页面中是否存在数据
	if roomStatus.App.InitialState.RoomStore.RoomInfo.Anchor.Nickname == "" {
		return base.GenAnchorInfoWhenErr(a.Anchor, fmt.Sprintf("https://live.douyin.com/%s", a.UID)),
			fmt.Errorf("不存在的直播间")
	}

	roomInfo := roomStatus.App.InitialState.RoomStore.RoomInfo
	anchorInfo := base.AnchorInfo{
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
	return Name
}

// GetStreamHeaders 请求直播流时的请求头
func (a *AnchorDouyin) GetStreamHeaders() map[string]string {
	return map[string]string{
		"user-agent": request.UAWin,
	}
}

// 提取网页中的 RENDER_DATA
func parseRenderData[T any](dyUrl string) (*T, error) {
	// 获取抖音网页文本
	text, err := request.Client.GetText(dyUrl, headers)
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
