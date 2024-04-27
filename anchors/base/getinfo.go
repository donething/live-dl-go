package base

import (
	"fmt"
	"time"
)

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
