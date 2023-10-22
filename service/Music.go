package service

import (
	"context"
	"encoding/base64"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/ioutil"
)

type Music struct {
	ctx context.Context
}

func NewMusic() *Music {
	return &Music{}
}

func (m *Music) Startup(ctx context.Context) {
	m.ctx = ctx
	println("音乐模块加载完毕")
}

// SelectMusic 选择音乐目录
func (m *Music) SelectMusic() string {
	base64Encoded := ""
	dialog, err := runtime.OpenFileDialog(m.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{{DisplayName: "音乐文件", Pattern: "*.mp3"}},
	})

	if err != nil {
		return ""
	}
	if dialog != "" {
		data, err := ioutil.ReadFile(dialog)
		if err != nil {
			return ""
		}
		base64Encoded = base64.StdEncoding.EncodeToString(data)
	}
	return base64Encoded
}
