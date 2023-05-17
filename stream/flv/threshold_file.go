package flv

import (
	"fmt"
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/utils-go/dofile"
	"io"
	"os"
)

// ThresholdFile 可限制文件大小的文件写入器
type ThresholdFile struct {
	// 视频输入流，将根据 stream 中的`CreateReader`自动创建
	reader io.ReadCloser
	// 保存的路径。会根据该路径生成临时路径作为实际的保存路径，避免文件重复
	path string
	// 限制当前文件的最大字节数
	threshold int64
	// 流的信息
	stream *Stream

	// 根据 path 自动生成的临时路径，避免重复
	uniPath string
	// 当前写入的文件
	file *os.File
}

// NewThresholdFile 创建实例
//
// 当只调用`Write()`时，`reader`可以传`nil`
func NewThresholdFile(path string, threshold int64, stream *Stream) *ThresholdFile {
	return &ThresholdFile{
		path:      path,
		threshold: threshold,
		stream:    stream,
	}
}

// StartSave 读取视频流，保存到可限制文件大小的文件中
func (f *ThresholdFile) StartSave() error {
	// 最后合并、发送该文件夹中的视频片段
	defer func() {
		f.file.Close()
		f.reader.Close()

		if exists, err := dofile.Exists(f.uniPath); err != nil || !exists {
			return
		}
		hanlders.ChHandle <- &hanlders.InfoHandle{
			Path:    f.uniPath,
			Title:   f.stream.Title,
			Handler: f.stream.Handler,
		}
	}()

	// 缓存
	var buf = make([]byte, 32*1024)
	for {
		// 是否停止录制视频流
		if f.stream.Stop.GetStop() {
			return nil
		}

		// 首次或换新文件存储时，要重新创建数据输入流
		if f.reader == nil {
			reader, err := f.stream.CreateReader()
			if err != nil {
				return fmt.Errorf("创建 flv 视频输入流出错：%w", err)
			}
			f.reader = reader
		}

		n, err := f.reader.Read(buf)
		// 读取出错
		if n < 0 {
			return fmt.Errorf("读取视频内容出错：%w", err)
		}
		// 已读完当前视频切片
		if n == 0 {
			break
		}

		// 写入数据
		_, err = f.Write(buf[:n])
		if err != nil {
			return fmt.Errorf("写入可限制大小的视频文件出错：%w", err)
		}
	}

	return nil
}

// Write 写入
func (f *ThresholdFile) Write(bs []byte) (int, error) {
	// 初始化文件
	if f.file == nil {
		// 打开写入的文件
		f.uniPath = dofile.UniquePath(f.path)
		file, err := os.Create(f.uniPath)
		if err != nil {
			return 0, fmt.Errorf("创建视频文件出错：%w", err)
		}
		f.file = file
	}

	// 写入
	n, err := f.file.Write(bs)
	if err != nil {
		return 0, fmt.Errorf("写入视频文件出错：%w", err)
	}
	f.stream.CurBytes.AddBytes(int64(n))

	// 判断是否需要更换新文件保存
	if f.threshold != 0 && f.stream.CurBytes.GetBytes() >= f.threshold {
		// 先关闭当前文件
		f.file.Close()

		// 再处理当前视频文件
		hanlders.ChHandle <- &hanlders.InfoHandle{
			Path:    f.uniPath,
			Title:   f.stream.Title,
			Handler: f.stream.Handler,
		}

		// 最后清空该视频文件的信息，以便新创建
		f.file = nil
		f.uniPath = ""
		f.stream.CurBytes.ResetBytes()
		f.reader.Close()
		f.reader = nil
	}

	return n, nil
}
