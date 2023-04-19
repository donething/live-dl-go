package entity

// IStream 接口
type IStream interface {
	// Start 开始录制
	Start() error

	GetStream() *Stream
}
