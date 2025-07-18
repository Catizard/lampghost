package main

import (
	"context"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/controller"
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/dto"
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
	*controller.CustomDiffTableController
	*controller.CustomCourseController
	*service.MonitorService
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
	conf, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	// auto-reload module
	monitorService, notifySyncChan := service.NewMonitorService(conf)

	// rival module
	rivalInfoService := service.NewRivalInfoService(db, monitorService, notifySyncChan)
	rivalTagService := service.NewRivalTagService(db)
	rivalScoreLogService := service.NewRivalScoreLogService(db)
	rivalSongDataService := service.NewRivalSongDataService(db)
	rivalScoreDataLogService := service.NewRivalScoreDataLogService(db)
	rivalInfoController := controller.NewRivalInfoController(rivalInfoService)
	rivalTagController := controller.NewRivalTagController(rivalTagService)
	rivalScoreLogController := controller.NewRivalScoreLogController(rivalScoreLogService)
	rivalScoreDataLogController := controller.NewRivalScoreDataLogController(rivalScoreDataLogService)
	rivalSongDataController := controller.NewRivalSongDataController(rivalSongDataService)

	// Set up the initial scorelog path
	if mainUser, err := rivalInfoService.QueryMainUser(); err == nil && mainUser != nil {
		monitorService.SetScoreLogFilePath(*mainUser.ScoreLogPath)
	}

	// download task module
	downloadTaskService := service.NewDownloadTaskService(db, conf, configService.Subscribe())
	downloadTaskController := controller.NewDownloadTaskController(downloadTaskService)

	// difficult table module
	diffTableService := service.NewDiffTableService(db, downloadTaskService)
	diffTableController := controller.NewDiffTableController(diffTableService)
	courseInfoService := service.NewCourseInfoSerivce(db, conf, configService.Subscribe())
	courseInfoController := controller.NewCourseInfoController(courseInfoService)

	// custom difficult table module
	customDiffTableService := service.NewCustomDiffTableService(db)
	customDiffTableController := controller.NewCustomDiffTableController(customDiffTableService)
	folderService := service.NewFolderService(db)
	folderController := controller.NewFolderController(folderService)
	customCourseService := service.NewCustomCourseService(db)
	customCourseController := controller.NewCustomCourseController(customCourseService)

	// Internal Server
	internalServer := server.NewInternalServer(
		customDiffTableService,
		customCourseService,
		folderService,
		rivalInfoService,
		rivalScoreLogService,
		rivalTagService,
		rivalSongDataService,
	)

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
		customDiffTableController,
		customCourseController,
		monitorService,
		internalServer,
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

func (a *App) GENERATOR_NOTIFICATION_DTO() *dto.NotificationDto { return nil }
