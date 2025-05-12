package main

import (
	"context"
	"log"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/controller"
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/server"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	*controller.ConfigController
	*controller.RivalInfoController
	*controller.RivalTagController
	*controller.RivalScoreLogController
	*controller.RivalScoreDataLogController
	*controller.RivalSongDataController
	*controller.DiffTableController
	*controller.CourseInfoController
	*controller.FolderController
	*server.InternalServer
}

// NewApp creates a new App application struct
func NewApp() *App {
	dbConfig := config.NewDatabaseConfig()
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("initialize database: %s", err)
	}

	// config module
	configService := service.NewConfigService(db)
	configController := controller.NewConfigController(configService)

	// rival module
	rivalInfoService := service.NewRivalInfoService(db)
	rivalTagService := service.NewRivalTagService(db)
	rivalScoreLogService := service.NewRivalScoreLogService(db)
	rivalSongDataService := service.NewRivalSongDataService(db)
	rivalScoreDataLogService := service.NewRivalScoreDataLogService(db)
	rivalInfoController := controller.NewRivalInfoController(rivalInfoService)
	rivalTagController := controller.NewRivalTagController(rivalTagService)
	rivalScoreLogController := controller.NewRivalScoreLogController(rivalScoreLogService)
	rivalScoreDataLogController := controller.NewRivalScoreDataLogController(rivalScoreDataLogService)
	rivalSongDataController := controller.NewRivalSongDataController(rivalSongDataService)

	// difficult table module
	diffTableService := service.NewDiffTableService(db)
	diffTableController := controller.NewDiffTableController(diffTableService)
	courseInfoService := service.NewCourseInfoSerivce(db)
	courseInfoController := controller.NewCourseInfoController(courseInfoService)

	// folder module
	folderService := service.NewFolderService(db)
	folderController := controller.NewFolderController(folderService)
	folderInternalServer := server.NewInternalServer(folderService)

	return &App{
		nil,
		configController,
		rivalInfoController,
		rivalTagController,
		rivalScoreLogController,
		rivalScoreDataLogController,
		rivalSongDataController,
		diffTableController,
		courseInfoController,
		folderController,
		folderInternalServer,
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Start internal server
	if err := a.RunServer(); err != nil {
		log.Fatalf("cannot start internal server: %s", err)
	}
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) OpenFileDialog(title string) result.RtnData {
	fp, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	if err != nil {
		return result.NewErrorData(err)
	}
	return result.NewRtnData(fp)
}
