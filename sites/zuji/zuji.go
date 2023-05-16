package zuji

import (
	"encoding/json"
	"fmt"
	"github.com/donething/live-dl-go/comm"
	"github.com/donething/live-dl-go/sites/entity"
	"regexp"
)

const (
	// 目标域名
	host = "szsbtech.com"

	Plat = "zuji"
	name = "足迹"
)

// AnchorZuji 足迹主播
type AnchorZuji struct {
	// 主播的 ID 为用户 ID
	*entity.Anchor
}

// 1. 获取 sessionid
func getSessionid(uid string) (string, error) {
	var headers = map[string]string{
		"User-Agent": comm.UAAndroid,
	}
	text, err := comm.Client.GetText(fmt.Sprintf("http://share-g3g5zb3o.i.%s/r/%s", host, uid), headers)
	if err != nil {
		return "", fmt.Errorf("执行获取 session 的请求出错：%w", err)
	}

	re := regexp.MustCompile(`sessionid\s*=\s*'(\w+?)'`)
	match := re.FindStringSubmatch(text)
	if len(match) < 2 {
		return "", fmt.Errorf("无法匹配到 sessionid")
	}

	return match[1], nil
}

// 2. 获取基础信息（除直播流地址以外的）
func getBasicInfo(uid string, sessionid string) (*RespInterface, error) {
	postHeaders := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"User-Agent":   comm.UAAndroid,
	}

	u := fmt.Sprintf("http://share-g3g5zb3o.i.%s/call_interface.php", host)
	data := fmt.Sprintf("joinroom=joinroom&room=%s&sessionid=%s", uid, sessionid)
	bs, err := comm.Client.PostForm(u, data, postHeaders)
	if err != nil {
		return nil, fmt.Errorf("执行获取主播基础信息的请求出错：%w", err)
	}

	var resp RespInterface
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, fmt.Errorf("解析主播基础信息数据时出错：%w", err)
	}

	if resp.Retval != "ok" {
		return nil, fmt.Errorf("获取主播基础信息数据失败：'%s %s'", resp.Retval, resp.Reterr)
	}

	return &resp, nil
}

// GetAnchorInfo 3. 获取足迹主播的信息
func (a *AnchorZuji) GetAnchorInfo() (*entity.AnchorInfo, error) {
	// 先获取基础信息
	sessionid, err := getSessionid(a.ID)
	if err != nil {
		return entity.GenAnchorInfoWhenErr(a.Anchor,
			fmt.Sprintf("https://share-aq2g4taz.i.%s/u/%s", host, a.ID)), err
	}

	vData, err := getBasicInfo(a.ID, sessionid)
	if err != nil {
		return entity.GenAnchorInfoWhenErr(a.Anchor,
			fmt.Sprintf("https://share-aq2g4taz.i.%s/u/%s", host, a.ID)), err
	}

	info := vData.Retinfo
	anchorInfo := entity.AnchorInfo{
		Anchor: a.Anchor,
		Avatar: info.Logourl,
		Name:   info.Nickname,
		WebUrl: fmt.Sprintf("http://share-g3g5zb3o.i.%s/r/%s", host, a.ID),
		Title:  info.Title,
		IsLive: info.Roomstatus == 1,
		Denied: info.Permission != 0,
	}

	// 如果主播不在播，就不用获取直播流地址了，直接返回已获取的信息
	if !anchorInfo.IsLive {
		return &anchorInfo, nil
	}

	// 主播在播，获取直播流地址
	// 使用 `Android` 端的 `User-Agent` 可以返回 `.m3u8` 流，`Windows` 端则返回 `rtmp` 流
	var headers = map[string]string{
		"User-Agent": comm.UAAndroid,
	}
	u := fmt.Sprintf("https://m.%s/appgw/v2/watchstart?vid=%s&sessionid=%s", host, info.Vid, sessionid)
	bs, err := comm.Client.GetBytes(u, headers)
	if err != nil {
		return nil, fmt.Errorf("执行获取直播流地址的请求出错：%w", err)
	}

	var resp RespWatchStart
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, fmt.Errorf("解析直播流地址的数据时出错：%w", err)
	}

	if resp.Retval != "ok" {
		return nil, fmt.Errorf("获取直播流地址失败：'%s %s'", resp.Retval, resp.Reterr)
	}

	anchorInfo.StreamUrl = resp.Retinfo.PlayURL

	return &anchorInfo, nil
}

// GetPlatName 获取平台名
func (a *AnchorZuji) GetPlatName() string {
	return name
}

// GetStreamHeaders 请求直播流时的请求头
func (a *AnchorZuji) GetStreamHeaders() map[string]string {
	return map[string]string{
		"user-agent": comm.UAAndroid,
	}
}
