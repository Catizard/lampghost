package controller

import (
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/charmbracelet/log"
)

type FolderController struct {
	folderService *service.FolderService
}

func NewFolderController(folderService *service.FolderService) *FolderController {
	return &FolderController{
		folderService: folderService,
	}
}

func (ctl *FolderController) AddFolder(name string) result.RtnMessage {
	log.Info("[Controller] Calling FolderController.AddFolder")
	_, err := ctl.folderService.AddFolder(name)
	if err != nil {
		log.Errorf("[FolderController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *FolderController) FindFolderTree() result.RtnDataList {
	log.Info("[Controller] Calling FolderController.FindFolderTree")
	rows, _, err := ctl.folderService.FindFolderTree()
	if err != nil {
		log.Errorf("[FolderController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}
