package main

import (
	"embed"

	"github.com/charmbracelet/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	log.SetLevel(log.DebugLevel)

	// Create an instance of the app structure
	app := NewApp()

	// Bindings
	var bind []any
	bind = append(bind, app)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "lampghost_wails",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind:             bind,
		Mac: &mac.Options{
			WebviewIsTransparent: true,
		},
		LogLevel: logger.DEBUG,
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
