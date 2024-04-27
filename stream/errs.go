package stream

import "errors"

var (
	ErrOnlivePermissionDenied = errors.New("缺少观看权限")
	ErrOnliveUnknownStream    = errors.New("未知的直播流类型")
)
