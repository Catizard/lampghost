package main

import (
	"context"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/controller"
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/server"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/charmbracelet/log"
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
	*controller.DownloadTaskController
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
	// NOTE: Some modules in lampghost are working based on config module. And obviously, we don't
	// want to query config current value each time. Therefore, config module & App (God object)
	// implements a notify mechanism:
	//
	// config module -> App: config has changed!
	// App -> DownloadTaskService: config has changed!
	// App -> ...Service: config has changed!
	configPublishChannel := make(chan any) // The config notify publish side
	configSubscribeChannel := make([]chan any, 0)
	configSubscribeChannel = append(configSubscribeChannel, make(chan any)) // DownloadTaskService
	// Here goes other subscribe channels
	go func() {
		for {
			notify := <-configPublishChannel
			// This should never stuck
			for _, ch := range configSubscribeChannel {
				ch <- notify
			}
		}
	}()
	configService := service.NewConfigService(db, configPublishChannel)
	configController := controller.NewConfigController(configService)
	conf, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

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

	// download task module
	downloadTaskService := service.NewDownloadTaskService(db, conf, configSubscribeChannel[0])
	downloadTaskController := controller.NewDownloadTaskController(downloadTaskService)

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
		downloadTaskController,
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
	a.RivalInfoController.InjectContext(ctx)
	a.DownloadTaskController.InjectContext(ctx)
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

func (a *App) OpenDirectoryDialog(title string) result.RtnData {
	fp, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	if err != nil {
		return result.NewErrorData(err)
	}
	return result.NewRtnData(fp)
}
