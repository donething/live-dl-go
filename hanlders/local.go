package hanlders

import (
	"fmt"
	"os"
	"path/filepath"
)

// LocalHanlder 保存到本地
type LocalHanlder struct{}

func (l *LocalHanlder) Handle(info *InfoHandle) error {
	dst := filepath.Join(filepath.Dir(info.Path), fmt.Sprintf("%s%s", info.Title, filepath.Ext(info.Path)))
	fmt.Printf("重命名：%s => %s\n", info.Path, dst)
	return os.Rename(info.Path, dst)
}
