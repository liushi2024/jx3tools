package main

import (
	"changeme/service"
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	keyboard := service.NewKeyboard()
	hotKey := service.NewHotKey()
	// Create application with options
	err := wails.Run(&options.App{
		Title:         "jx3tool",
		Width:         330,
		Height:        600,
		DisableResize: true,
		Frameless:     true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			keyboard.Startup(ctx)
			hotKey.Startup(keyboard)
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
