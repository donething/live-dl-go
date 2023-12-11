package entity

import (
	"fmt"
	"time"
)

// Anchor 需要获取直播流的主播的信息
type Anchor struct {
	// 直播流
	UID string `json:"uid"`
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

	// 是否设置了权限，无法观看
	Denied bool `json:"hasPermission"`
}

const (
	// MaxRetry 获取主播信息失败后，重试的次数
	MaxRetry = 3
)

// GenAnchorInfoWhenErr 当 GetAnchorInfo() 获取主播信息出错时，避免返回 nil，而是快速生成实例
func GenAnchorInfoWhenErr(anchor *Anchor, webUrl string) *AnchorInfo {
	return &AnchorInfo{
		Anchor: anchor,
		Name:   fmt.Sprintf("%s %s", anchor.Plat, anchor.UID),
		Title:  "获取出错",
		WebUrl: webUrl,
	}
}

// TryGetAnchorInfo 获取主播信息，可指定失败后的重试次数
func TryGetAnchorInfo(anchorSite IAnchor, retry int) (*AnchorInfo, error) {
	fail := 0
	var info *AnchorInfo
	var err error

	for {
		info, err = anchorSite.GetAnchorInfo()
		if err != nil {
			// 重试
			if fail < retry {
				fail++
				time.Sleep(1 * time.Second)
				continue
			}

			return nil, err
		}

		// 获取成功
		return info, nil
	}
}
