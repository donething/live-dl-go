package capture_status

import (
	"fmt"
	"github.com/donething/live-dl-go/anchors/base"
	"sync"
)

// CapStatus 正在录制的主播的状态
type CapStatus[T any] struct {
	cap map[string]T
	mu  sync.Mutex
}

// New 生成 CapStatus 的实例
func New[T any]() *CapStatus[T] {
	return &CapStatus[T]{
		cap: make(map[string]T),
	}
}

// Set 设置
func (c *CapStatus[T]) Set(k string, v T) {
	c.mu.Lock()
	c.cap[k] = v
	c.mu.Unlock()
}

// Get 读取
func (c *CapStatus[T]) Get(k string) (T, bool) {
	var v T
	var ok bool

	c.mu.Lock()
	v, ok = c.cap[k]
	c.mu.Unlock()

	return v, ok
}

// Del 删除
func (c *CapStatus[T]) Del(k string) {
	c.mu.Lock()
	delete(c.cap, k)
	c.mu.Unlock()
}

// Keys 返回所有键
func (c *CapStatus[T]) Keys() []string {
	keys := make([]string, 0, len(c.cap))
	c.mu.Lock()
	for key := range c.cap {
		keys = append(keys, key)
	}
	c.mu.Unlock()

	return keys
}

// GenCapturingKey 正在录制的主播的键，避免重复录制。格式如 "<平台>_<主播ID>"，如 "bili_12345"
func GenCapturingKey(anchor base.Anchor) string {
	return fmt.Sprintf("%s_%s", anchor.Plat, anchor.UID)
}
