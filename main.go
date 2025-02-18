package main

import (
	"embed"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/controller"
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/server"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/charmbracelet/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
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
	if err := c.Provide(service.NewRivalInfoService); err != nil {
		panic(err)
	}
	if err := c.Provide(service.NewDiffTableService); err != nil {
		panic(err)
	}
	if err := c.Provide(service.NewRivalTagService); err != nil {
		panic(err)
	}
	if err := c.Provide(service.NewRivalSongDataService); err != nil {
		panic(err)
	}
	if err := c.Provide(service.NewFolderService); err != nil {
		panic(err)
	}
	if err := c.Provide(service.NewCourseInfoSerivce); err != nil {
		panic(err)
	}
	if err := c.Provide(service.NewRivalScoreLogService); err != nil {
		panic(err)
	}
	if err := c.Provide(controller.NewRivalInfoController); err != nil {
		panic(err)
	}
	if err := c.Provide(controller.NewDiffTableController); err != nil {
		panic(err)
	}
	if err := c.Provide(controller.NewFolderController); err != nil {
		panic(err)
	}
	if err := c.Provide(controller.NewRivalScoreLogController); err != nil {
		panic(err)
	}
	if err := c.Provide(controller.NewRivalSongDataController); err != nil {
		panic(err)
	}
	if err := c.Provide(controller.NewCourseInfoController); err != nil {
		panic(err)
	}
	if err := c.Provide(server.NewInternalServer); err != nil {
		panic(err)
	}

	var bind []interface{}
	if err := c.Invoke(func(controller *controller.RivalInfoController) error {
		bind = append(bind, controller)
		return nil
	}); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(controller *controller.DiffTableController) error {
		bind = append(bind, controller)
		return nil
	}); err != nil {
		panic(err)
	}

	if err := c.Invoke(func(controller *controller.FolderController) error {
		bind = append(bind, controller)
		return nil
	}); err != nil {
		panic(err)
	}

	if err := c.Invoke(func(controller *controller.RivalScoreLogController) error {
		bind = append(bind, controller)
		return nil
	}); err != nil {
		panic(err)
	}

	if err := c.Invoke(func(controller *controller.RivalSongDataController) error {
		bind = append(bind, controller)
		return nil
	}); err != nil {
		panic(err)
	}

	if err := c.Invoke(func(controller *controller.CourseInfoController) error {
		bind = append(bind, controller)
		return nil
	}); err != nil {
		panic(err)
	}

	if err := c.Invoke(func(server *server.InternalServer) error {
		return server.RunServer()
	}); err != nil {
		panic(err)
	}

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
