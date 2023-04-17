package entity

import "github.com/donething/live-dl-go/hanlders"

// IStream 接口，用于
type IStream interface {
	// Start 开始录制
	Start() error

	GetChErr() chan error

	GetChRestart() chan bool

	// Reset Stream 中需要设为指定值的参数，设为默认值的不需传递参数
	Reset(title, streamUrl string, headers map[string]string, path string,
		fileSizeThreshold int, hanlder hanlders.IHandler)
}
