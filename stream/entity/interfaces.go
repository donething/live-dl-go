package entity

// IStream 接口
type IStream interface {
	GetStream() *Stream

	Capture() error
}
