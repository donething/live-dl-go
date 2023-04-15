package logger

import "github.com/donething/utils-go/dolog"

// 输出日志
var Info, Warn, Error = dolog.InitLog(dolog.DefaultFlag)
