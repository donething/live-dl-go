package comm

import "github.com/donething/utils-go/dohttp"

const (
	// SizeOneGB 1 GB 字节
	SizeOneGB = 1024 * 1024 * 1024
)

const (
	UAWin = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"

	UAAndroid = "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/112.0.5615.48 Mobile Safari/537.36"
)

var Client = dohttp.New(false, false)
