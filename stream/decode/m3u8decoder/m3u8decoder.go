// Package m3u8decoder 解析 m3u8 文件
//
// 实例化 M3u8Decoder 后，调用 Decode 解析
package m3u8decoder

import (
	"bufio"
	"bytes"
	"fmt"
	dllive "github.com/donething/live-dl-go/request"
	"io"
	"strings"
)

// M3u8Segment 包含视频切片的信息
type M3u8Segment struct {
	URL string
}

// M3u8Decoder 解码m3u8文件
type M3u8Decoder struct {
	Segments []*M3u8Segment
}

// New 创建实例
func New() *M3u8Decoder {
	return &M3u8Decoder{}
}

// Decode 从Reader中解码m3u8文件
func (d *M3u8Decoder) Decode(url string, headers map[string]string) error {
	// 读取m3u8文件的内容
	resp, err := dllive.Client.Get(url, headers)
	if err != nil {
		return fmt.Errorf("请求m3u8文件出错：%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return fmt.Errorf("读取m3u8文件的响应码：%s (URL: %s)", resp.Status, url)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取m3u8文件出错：%w", err)
	}

	if !strings.Contains(string(bs), "#EXTM3U") {
		return fmt.Errorf("请求的文件不是m3u8格式")
	}

	// 逐行读取
	scanner := bufio.NewScanner(bytes.NewReader(bs))
	// 需要在切片URL的前面加地址（"http://example.com/dir"，也可为""）
	var prefix *string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		segment := new(M3u8Segment)

		// 还没有判断是否需要添加前缀，开始先判断
		if prefix == nil {
			prefix = new(string)
			*prefix = getPrefix(url, line)
		}

		segment.URL = *prefix + line
		d.Segments = append(d.Segments, segment)
		segment = nil
	}
	return nil
}

// 返回完整的切片地址
func getPrefix(m3u8Url, segmentUrl string) string {
	// 切片地址为 "https?://"，不需要添加前缀。直接返回 ""
	if strings.HasPrefix(segmentUrl, "http") {
		return ""
	}

	// 切片地址为 "/record/1.ts"，需要追加`域名` "https://example.com"
	if strings.HasPrefix(segmentUrl, "/") {
		// 得到域名开始的索引，即后面的字符串为 "host.com/path/..."
		hostStartIndex := strings.Index(m3u8Url, "//") + 2
		// 得到从域名开始（即去除协议后的 "host.com/path/..."）的后的第一个'/'的索引
		pathStartIndexFromHost := strings.Index(m3u8Url[hostStartIndex:], "/")

		return m3u8Url[:hostStartIndex+pathStartIndexFromHost]
	}

	// 切片地址为 "1.ts"，需要追加`域名+目录路径` "https://example.com/record/" 后返回
	return m3u8Url[:strings.LastIndex(m3u8Url, "/")+1]
}
