package files

import (
	"fmt"
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/live-dl-go/stream/entity"
	"github.com/donething/utils-go/dofile"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ThresholdFile 可限制文件大小的文件写入器
type ThresholdFile struct {
	// 数据输入流
	reader io.ReadCloser
	// 保存的路径。会根据该路径生成临时路径作为实际的保存路径，避免文件重复
	path string
	// 限制当前文件的最大字节数
	threshold int64
	// 流的信息
	stream *entity.Stream
	// 新文件是否需要重读数据输入流。flv 视频需要为 true 来获取新文件的视频头信息；m3u8 视频流可根据流的失效时间判断是否传 ture
	needRecreate bool

	// 根据 path 自动生成的临时路径，避免重复
	uniPath string
	// 当前写入的文件
	file *os.File
	// 当前文件已写入的字节数
	bytes int64
}

// NewThresholdFile 创建实例
//
// 当只调用`Write()`时，`reader`可以传`nil`
func NewThresholdFile(reader io.ReadCloser, needRecreate bool, path string, threshold int64,
	stream *entity.Stream) *ThresholdFile {
	return &ThresholdFile{
		reader:       reader,
		needRecreate: needRecreate,
		path:         path,
		threshold:    threshold,
		stream:       stream,
	}
}

// StartSave 读取视频流，保存到可限制文件大小的文件中
func (f *ThresholdFile) StartSave() error {
	defer f.reader.Close()
	defer f.file.Close()
	// 最后合并、发送该文件夹中的视频片段
	defer func() {
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
		f.uniPath = filepath.Join(filepath.Dir(f.path),
			fmt.Sprintf("%d_%s", time.Now().UnixMilli(), filepath.Base(f.path)))
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
	f.bytes += int64(n)

	// 判断是否需要更换新文件保存
	if f.threshold != 0 && f.bytes >= f.threshold {
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
		f.bytes = 0

		// 换新文件存储时，可能要重新创建数据输入流
		if f.needRecreate {
			f.reader.Close()
			reader, err := f.stream.CreateReader()
			if err != nil {
				return 0, fmt.Errorf("重新创建数据输入流出错：%w", err)
			}
			f.reader = reader
		}
	}

	return n, nil
}
