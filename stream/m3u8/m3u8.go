package m3u8

import (
	"fmt"
	"github.com/donething/live-dl-go/hanlders"
	"github.com/donething/live-dl-go/request"
	"github.com/donething/live-dl-go/stream/basestream"
	"github.com/donething/live-dl-go/stream/decode/m3u8decoder"
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dovideo"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Stream m3u8 直播流
type Stream struct {
	*basestream.Stream
}

// NewStream 创建 Stream 的实例
//
// 参数 path 视频的保存路径，以 ".mp4" 结尾
func NewStream(task *hanlders.TaskInfo, streamUrl string, headers map[string]string) basestream.IStream {
	return &Stream{
		Stream: &basestream.Stream{
			TaskInfo:  task,
			StreamUrl: streamUrl,
			Headers:   headers,
		},
	}
}

func (s *Stream) GetStream() *basestream.Stream {
	return s.Stream
}

// Capture 录制 m3u8 视频流
func (s *Stream) Capture() error {
	// 避免重复 URL
	urlsHistory := NewUrlsHistory(MaxUrlsHistory)
	// 临时文件的基目录
	baseDir := filepath.Dir(s.Path)
	// 当前写入的文件夹
	folder := genFolderName(s.Path)
	// 偶尔获取 m3u8 文件会失败，尝试重试
	retry := 0

	// 最后合并、发送该文件夹中的视频片段
	defer func() {
		p := filepath.Join(baseDir, folder)
		if exists, err := dofile.Exists(p); err != nil || !exists {
			return
		}

		err := concatAndSend(p, s)
		if err != nil {
			fmt.Printf("合并、发送最后的文件夹的视频出错：%s\n", err.Error())
		}
	}()

	for {
		// 是否停止录制视频流。放在此层 for，避免频繁判断
		if s.Stop.GetStop() {
			return nil
		}

		// 解码 m3u8 视频列表
		m := m3u8decoder.New()
		err := m.Decode(s.StreamUrl, s.Headers)
		if err != nil {
			// 获取失败时可以重试2次
			if retry < 2 {
				retry++
				time.Sleep(1 * time.Second)
				continue
			}
			return fmt.Errorf("解码 m3u8 文件出错：%w", err)
		}

		// 如果没有获取到新的切片，表示直播结束
		if len(m.Segments) == 0 {
			break
		}

		// 创建临时目录
		err = os.MkdirAll(filepath.Join(baseDir, folder), 0755)
		if err != nil {
			return fmt.Errorf("创建 m3u8 视频保存目录出错：%w", err)
		}

		// 下载视频片段
		for _, seg := range m.Segments {
			if exists := urlsHistory.Exists(seg.URL); exists {
				continue
			}

			resp, err := request.Client.Get(seg.URL, s.Headers)
			if err != nil {
				return fmt.Errorf("创建 m3u8 视频输入流出错。请求视频出错：%w", err)
			}
			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				return fmt.Errorf("创建 m3u8 视频输入流出错。读取视频的响应码：%s (URL: %s)", resp.Status, seg.URL)
			}

			curPath := filepath.Join(baseDir, folder, fmt.Sprintf("%d.ts", time.Now().UnixMilli()))
			file, err := os.Create(curPath)
			if err != nil {
				return fmt.Errorf("创建 m3u8 视频片段的文件出错：%w", err)
			}
			n, err := io.Copy(file, resp.Body)
			if err != nil {
				return fmt.Errorf("写入 m3u8 视频片段的文件出错：%w", err)
			}
			resp.Body.Close()
			file.Close()

			s.CurBytes.AddBytes(n)
		}

		// 在当前文件夹写入的数据已达到限制的大小，将视频保存到新文件夹中
		// 为了减少调用`os.MkdirAll`，放在循环下载切片的父层的此处
		if s.FileSizeThreshold != 0 && s.CurBytes.GetBytes() >= s.FileSizeThreshold {
			// 合并、发送该文件夹中的视频片段
			err = concatAndSend(filepath.Join(baseDir, folder), s)
			if err != nil {
				return fmt.Errorf("合并、发送当前文件夹的视频出错：%w", err)
			}

			// 新的视频片段文件夹
			folder = genFolderName(s.Path)
			// 注意：清除当前文件夹的数据记录，以便重新计算下一个文件夹的数据
			s.CurBytes.ResetBytes()
		}

		// time.Sleep(time.Duration(domath.RandInt(3, 6)))
	}

	return nil
}

// 合并、发送视频
func concatAndSend(dir string, s *Stream) error {
	// 合并视频
	unique := dofile.UniquePath(s.Path)
	err := dovideo.Concat(dir, ".ts", unique)
	if err != nil {
		return fmt.Errorf("合并目录中的视频出错(%s)：\n%w", dir, err)
	}

	// 删除临时目录
	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("合并视频时，删除临时目录出错(%s)：%w", dir, err)
	}

	// 再处理当前合并的视频文件
	s.TaskInfo.Path = unique
	hanlders.ChHandle <- s.TaskInfo

	return nil
}

// 生成录制当前视频流的工作文件夹
func genFolderName(path string) string {
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return fmt.Sprintf("%s_%d", name, time.Now().Unix())
}
