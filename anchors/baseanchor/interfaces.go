package baseanchor

// IAnchor 主播的接口
type IAnchor interface {
	// GetAnchorInfo 获取主播信息
	// 约定：有错误发生时，AnchorInfo 不要返回 nil
	// 此时用 GenAnchorInfoWhenErr() 快速生成实例
	GetAnchorInfo() (*AnchorInfo, error)

	// GetPlatName 获取平台名。如"哔哩哔哩"
	GetPlatName() string

	// GetStreamHeaders 请求直播流时的请求头
	GetStreamHeaders() map[string]string
}
