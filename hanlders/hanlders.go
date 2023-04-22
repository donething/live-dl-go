package hanlders

import (
	"github.com/donething/live-dl-go/comm/logger"
)

// IHandler 处理已完成下载的视频的处理器接口
type IHandler interface {
	Handle(info *InfoHandle) error
}

// InfoHandle 需要被处理的文件和信息
type InfoHandle struct {
	// 文件路径
	Path string
	// 标题。注意作为TG的caption时，需要转义
	Title string
	// 文件处理器。在 handler worker 中将调用，来处理文件
	Handler IHandler
	// 上传到TG时的分段大小，为 0 不分段
	FileSizeThreshold int64
	// 	上传成功是否保留源文件
	Reserve bool
}

// 处理文件的goroutine数量
const handlerCount = 5

// ChHandle 发送需要被处理的文件的通道
var ChHandle = make(chan *InfoHandle, handlerCount)

// WGHandler 等待处理完所有视频后，才能退出
// var WGHandler = sync.WaitGroup{}

func init() {
	// 启动 goroutine 来完成工作
	for gr := 1; gr <= handlerCount; gr++ {
		go handler()
	}
	logger.Info.Println("视频处理器的 goroutine 已准备就绪")
}

// 视频文件处理器。从 ChHandle 接收视频文件
func handler() {
	for {
		info, ok := <-ChHandle
		if !ok {
			break
		}

		// 	实际工作
		logger.Info.Printf("开始处理视频文件：%s\n", info.Path)
		err := info.Handler.Handle(info)
		if err != nil {
			logger.Error.Printf("处理视频文件出错(%s)：%s\n", info.Path, err)
			continue
		}
	}
}
