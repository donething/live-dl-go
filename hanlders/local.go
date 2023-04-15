package hanlders

import (
	"dl-live-go/comm/logger"
	"fmt"
	"os"
	"path/filepath"
)

// LocalHanlder 保存到本地
type LocalHanlder struct{}

func (l *LocalHanlder) Handle(info *InfoHandle) error {
	dst := filepath.Join(filepath.Dir(info.Path), fmt.Sprintf("%s%s", info.Title, filepath.Ext(info.Path)))
	logger.Info.Printf("重命名：%s => %s\n", info.Path, dst)
	return os.Rename(info.Path, dst)
}
