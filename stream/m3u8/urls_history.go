package m3u8

const MaxUrlsHistory = 10

// UrlsHistory 避免发送、下载重复的 URL 切片
//
// 用一个固定容量的切片保存历史 URL
type UrlsHistory struct {
	// 已记录的 URL
	urls []string
	// 容量
	max int
	// 已记录的 URL 的数量
	count int
}

// NewUrlsHistory 创建保存历史 URL 的实例。参数为容量，可以设为 MaxUrlsHistory
func NewUrlsHistory(max int) *UrlsHistory {
	urls := make([]string, max)
	return &UrlsHistory{
		urls: urls,
		max:  max,
	}
}

// Exists 是否存在重复的切片 URL
func (u *UrlsHistory) Exists(url string) bool {
	// 存在该 URL
	for _, s := range u.urls {
		if url == s {
			return true
		}
	}

	// 不存在时，记录该 URL
	// fmt.Printf("保存 %d 切片 URL：%s\n", u.count%u.max, url)
	u.urls[u.count%u.max] = url
	u.count++

	return false
}
