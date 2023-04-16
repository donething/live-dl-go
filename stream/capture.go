// Package stream 捕获视频流

package stream

import (
	"fmt"
	"github.com/donething/live-dl-go/comm"
	"github.com/donething/live-dl-go/comm/logger"
	"github.com/donething/live-dl-go/hanlders"
	"os"
	"path/filepath"
	"time"
)

// PrepareCapture 录制直播流到文件
func (s *Stream) PrepareCapture() error {
	s.Path = filepath.Join(filepath.Dir(s.Path),
		fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(s.Path)))
	// 打开写入的文件
	file, err := os.Create(s.Path)
	if err != nil {
		return fmt.Errorf("创建视频文件出错：%w", err)
	}

	go s.save(file)
	return nil
}

// 因为用`for url := range s.ChSegUrl`监听channel，需要主goroutine继续执行后面的语句来发送数据，
// 所以要在子goroutine中运行，以免死锁：`go s.save()`
func (s *Stream) save(file *os.File) {
	var restart = false
	// 需要穿参数，因为可能在等待 channel 传数据时，restart 为 true，然后`PrepareCapture()`重新执行
	// 导致 s.Path 路径改变，导致发送给 ChHandle 的 Path 错误
	defer func(p string) {
		// 传递是否需要重新开始直播流流保存到新文件
	LabelRestart:
		for {
			select {
			case s.ChRestart <- restart:
				break LabelRestart
			default:
				// 	继续等待
				logger.Info.Printf("等待发送新开始录制的信号…\n")
				time.Sleep(1 * time.Second)
			}
		}

		// 	后台处理视频
		file.Close()
	LabelHandle:
		for {
			select {
			case hanlders.ChHandle <- &hanlders.InfoHandle{
				Path:    p,
				Title:   s.Title,
				Handler: s.Handler,
			}:
				break LabelHandle
			default:
				// 	继续等待
				logger.Info.Printf("等待发送文件到视频处理器…(%s)\n", p)
				time.Sleep(1 * time.Second)
			}
		}
	}(s.Path)

	// 当前写入文件的字节计数器
	var fileSizeCounter = 0
	// 缓存
	var buf = make([]byte, 32*1024)

	// 获取视频（切片）的二进制数据
LabelSegs:
	for url := range s.ChSegUrl {
		// logger.Info.Printf("收到视频切片的链接：%s\n", url)
		resp, err := comm.Client.Get(url, s.Headers)
		if err != nil {
			s.ChErr <- fmt.Errorf("请求视频出错：%w", err)
			return
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			s.ChErr <- fmt.Errorf("读取视频的响应码：%s。URL: %s", resp.Status, url)
			return
		}

		for {
			n, err := resp.Body.Read(buf)
			// 读取出错
			if n < 0 {
				s.ChErr <- fmt.Errorf("读取视频内容出错：%w", err)
				return
			}
			// 已读完当前视频切片
			if n == 0 {
				break
			}

			// 写入数据
			_, err = file.Write(buf[:n])
			if err != nil {
				s.ChErr <- fmt.Errorf("写入视频文件出错：%w", err)
				return
			}

			// 统计当前文件的大小
			// 如果视频大小达到阈值，就重新开始保存直播流到新文件中
			// 需要重新打开直播流读取，这样避免手动为视频添加头信息
			if fileSizeCounter += n; s.FileSizeThreshold != 0 && fileSizeCounter >= s.FileSizeThreshold {
				fileSizeCounter = 0
				restart = true
				// 当下载m3u8时，s.ChSegUrl 一直有片段可以接收，仅通过 break 无法跳出2个循环
				break LabelSegs
			}
		}
	}

	// 没有错误时，需要传回 nil，避免阻塞
LabelErrNil:
	for {
		select {
		case s.ChErr <- nil:
			break LabelErrNil
		default:
			// 	继续等待
			logger.Info.Printf("等待发送无错信号…(%s)\n", s.Path)
			time.Sleep(1 * time.Second)
		}
	}
}
