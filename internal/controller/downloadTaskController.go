package controller

import (
	"context"

	"github.com/Catizard/lampghost_wails/internal/config/download"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/charmbracelet/log"
)

type DownloadTaskController struct {
	downloadTaskService *service.DownloadTaskService
}

func NewDownloadTaskController(downloadTaskService *service.DownloadTaskService) *DownloadTaskController {
	return &DownloadTaskController{
		downloadTaskService: downloadTaskService,
	}
}

func (ctl *DownloadTaskController) InjectContext(ctx context.Context) {
	ctl.downloadTaskService.InjectContext(ctx)
}

func (ctl *DownloadTaskController) SubmitDownloadTask(url string, taskName *string) result.RtnMessage {
	log.Info("[Controller] calling DownloadTaskController.SubmitDownloadTask")
	if err := ctl.downloadTaskService.SubmitDownloadTask(url, nil); err != nil {
		log.Errorf("[DownloadTaskController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *DownloadTaskController) SubmitSingleMD5DownloadTask(md5 string, taskName *string) result.RtnMessage {
	log.Info("[Controller] calling DownloadTaskController.SubmitSingleMD5DownloadTask")
	if err := ctl.downloadTaskService.SubmitSingleMD5DownloadTask(md5, taskName); err != nil {
		log.Errorf("[DownloadTaskController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *DownloadTaskController) FindDownloadTaskList() result.RtnDataList {
	log.Info("[Controller] calling DownloadTaskController.FindDownloadTaskList")
	rows, _, err := ctl.downloadTaskService.FindDownloadTaskList()
	if err != nil {
		log.Errorf("[DownloadTaskController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *DownloadTaskController) CancelDownloadTask(taskID uint) result.RtnMessage {
	log.Info("[Controller] calling DownloadTaskController.CancelDownloadTask")
	if err := ctl.downloadTaskService.CancelDownloadTask(taskID); err != nil {
		log.Errorf("[DownloadTaskController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *DownloadTaskController) RestartDownloadTask(taskID uint) result.RtnMessage {
	log.Info("[Controller] calling DownloadTaskController.RestartDownloadTask")
	if err := ctl.downloadTaskService.RestartDownloadTask(taskID); err != nil {
		log.Errorf("[DownloadTaskController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (*DownloadTaskController) GENERATOR_DOWNLOAD_TASK() *download.DownloadSource { return nil }
