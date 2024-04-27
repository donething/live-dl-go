package hanlders

// AfterHandleFunc 执行处理视频后的回调函数
type AfterHandleFunc func(task *TaskInfo, err error)

// IHandler 处理已完成下载的视频的处理器接口
type IHandler interface {
	Handle(task *TaskInfo)
}

// TaskInfo 需要被处理的文件和信息
type TaskInfo struct {
	Path              string          // 任务文件（保存）的路径。注意根据视频流的类型设置正确的文件格式，如".flv"、".mp4"等
	Title             string          // 标题。注意作为 TG 的 caption 时，需要转义
	Handler           IHandler        // 文件处理器。在 handler worker 中将调用，来处理文件
	FileSizeThreshold int64           // 文件的最大字节数，为 0 表示不分段。上传 TG 建议设为 hanlders.FileSizeThreshold
	AfterHandle       AfterHandleFunc // 执行完 Handler 后的回调
}

// 处理文件的 goroutine 数量
const handlerCount = 5

// ChHandle 发送需要被处理的文件的通道
var ChHandle = make(chan *TaskInfo, handlerCount)

// WGHandler 等待处理完所有视频后，才能退出
// var WGHandler = sync.WaitGroup{}

func init() {
	// 启动 goroutine 来完成工作
	for gr := 1; gr <= handlerCount; gr++ {
		go handler()
	}
}

// 视频文件处理器。从 ChHandle 接收视频文件
func handler() {
	for {
		task, ok := <-ChHandle
		if !ok {
			break
		}

		// 	实际工作
		task.Handler.Handle(task)
	}
}
