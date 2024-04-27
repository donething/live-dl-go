package basestream

// IStream 接口
type IStream interface {
	GetStream() *Stream

	Capture() error
}
