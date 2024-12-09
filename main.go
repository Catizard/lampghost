package main

import (
	"embed"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/charmbracelet/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"go.uber.org/dig"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	log.SetLevel(log.DebugLevel)
	c := dig.New()
	if err := c.Provide(config.NewDatabaseConfig); err != nil {
		panic(err)
	}
	if err := c.Provide(database.NewDatabase); err != nil {
		panic(err)
	}
	c.Invoke(func(db *database.DB) error {
		// Invoke database initialize
		return nil
	})
	// Create an instance of the app structure
	app := NewApp()

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
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			WebviewIsTransparent: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
