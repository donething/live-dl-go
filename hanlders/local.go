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

func (l *LocalHanlder) Handle(task *TaskInfo) {
	name := fmt.Sprintf("%s_%d%s", task.Title, time.Now().Unix(), filepath.Ext(task.Path))
	name = dofile.ValidFileName(name, "_")

	dst := filepath.Join(filepath.Dir(task.Path), name)
	// logger.Info.Printf("重命名：'%s' => '%s'\n", task.Path, dst)

	if task.AfterHandle != nil {
		task.AfterHandle(task, os.Rename(task.Path, dst))
	}
}
