package controller

import "github.com/Catizard/lampghost_wails/internal/service"

type FolderController struct {
	folderService *service.FolderService
}

func NewFolderController(folderService *service.FolderService) *FolderController {
	return &FolderController{
		folderService: folderService,
	}
}
