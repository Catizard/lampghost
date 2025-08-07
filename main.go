package main

import (
	"embed"
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/logger"
	"github.com/charmbracelet/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"golang.design/x/clipboard"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	log.SetLevel(log.DebugLevel)
	// Setup clipboard
	if err := clipboard.Init(); err != nil {
		log.Errorf("Cannot initialize clipboard due to %s, clipboard related feature may not work!", err)
	}

	// Create an instance of the app structure
	app := NewApp()

	// Bindings
	var bind []any
	bind = append(bind, app)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  fmt.Sprintf("lampghost_wails %s", config.VERSION),
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
		LogLevel: 3,
		Logger:   &logger.WailsLogger{},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
