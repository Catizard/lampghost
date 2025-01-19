package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
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

func (ctl *FolderController) DelFolder(ID uint) result.RtnMessage {
	log.Info("[Controller] calling FolderController.DelDiffTableHeader")
	err := ctl.folderService.DelFolder(ID)
	if err != nil {
		log.Errorf("[FolderController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *FolderController) DelFolderContent(contentID uint) result.RtnMessage {
	log.Info("[Controller] calling FolderController.DelFolderContent")
	err := ctl.folderService.DelFolderContent(contentID)
	if err != nil {
		log.Errorf("[FolderController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *FolderController) BindSongToFolder(diffTableDataID uint, folderIDs []uint) result.RtnMessage {
	log.Info("[Controller] Calling FolderController.AddSongToFolder")
	err := ctl.folderService.BindSongToFolder(diffTableDataID, folderIDs)
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

func (ctl *FolderController) FindFolderList() result.RtnDataList {
	log.Info("[Controller] Calling FolderController.FindFolderList")
	rows, _, err := ctl.folderService.FindFolderList()
	if err != nil {
		log.Errorf("[FolderController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *FolderController) SyncSongData() result.RtnMessage {
	panic("todo!")
}

func (ctl *FolderController) QueryFolderDefinition() result.RtnDataList {
	log.Info("[Controller] Calling FolderController.GenerateJson")
	rows, _, err := ctl.folderService.GenerateJson()
	if err != nil {
		log.Errorf("[FolderController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (FolderController) GENERATOR_FOLDER() *entity.Folder                          { return nil }
func (FolderController) GENERATOR_FOLDER_DTO() *dto.FolderDto                      { return nil }
func (FolderController) GENERATOR_FOLDER_CONTENT_DTO() *dto.FolderContentDto       { return nil }
func (FolderController) GENERATOR_FOLDER_DEFINITION_DTO() *dto.FolderDefinitionDto { return nil }
