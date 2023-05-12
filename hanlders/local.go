package hanlders

import (
	"fmt"
	"github.com/donething/utils-go/dofile"
	"os"
	"path/filepath"
	"time"
)

// LocalHanlder 保存到本地
type LocalHanlder struct{}

func (l *LocalHanlder) Handle(info *InfoHandle) error {
	name := fmt.Sprintf("%s_%d%s", info.Title, time.Now().UnixMilli(), filepath.Ext(info.Path))
	name = dofile.ValidFileName(name, "_")

	dst := filepath.Join(filepath.Dir(info.Path), name)
	// logger.Info.Printf("重命名：'%s' => '%s'\n", info.Path, dst)
	return os.Rename(info.Path, dst)
}
