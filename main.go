package main

import (
	"changeme/service"
	"context"
	"embed"
	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed public/icon.ico
var imageBytes []byte

//go:embed public/dd.dll
var dll []byte

func main() {
	app := NewApp()
	keyboard := service.NewKeyboard()
	hotKey := service.NewHotKey()
	// Create application with options
	err := wails.Run(&options.App{
		Title:         "按键小助手",
		Width:         330,
		Height:        600,
		DisableResize: true,
		Frameless:     true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			keyboard.Startup(ctx, dll)
			hotKey.Startup(keyboard)
			// 启动 systray
			systray.Run(func() {
				// 在任务栏上创建一个图标
				systray.SetIcon(imageBytes)
				systray.SetTitle("按键小助手")
				systray.SetTooltip("按键小助手")

				// 添加一个退出菜单项
				mQuit := systray.AddMenuItem("退出", "退出程序")
				mOpen := systray.AddMenuItem("显示", "显示界面")

				// 启动一个 goroutine 监听退出事件
				go func() {
					defer log.Println("菜单监听结束了")
					for {
						select {
						case <-mOpen.ClickedCh:
							runtime.WindowShow(ctx)
						case <-mQuit.ClickedCh:
							runtime.Quit(ctx)
							systray.Quit()
						}
					}
				}()
			}, func() {})
		},
		OnShutdown: func(ctx context.Context) {
			hotKey.Close()
			println("关闭了应用")
			//todo 记录窗口关闭
		},
		Bind: []interface{}{
			app,
			keyboard,
			hotKey,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
