package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
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

func (ctl *FolderController) FindFolderTree(filter *vo.FolderVo) result.RtnDataList {
	log.Info("[Controller] Calling FolderController.FindFolderTree")
	rows, _, err := ctl.folderService.FindFolderTree(filter)
	if err != nil {
		log.Errorf("[FolderController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *FolderController) FindFolderList(filter *vo.FolderVo) result.RtnDataList {
	log.Info("[Controller] Calling FolderController.FindFolderList")
	rows, _, err := ctl.folderService.FindFolderList(filter)
	if err != nil {
		log.Errorf("[FolderController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *FolderController) FindFolderContentList(filter *vo.FolderContentVo) result.RtnDataList {
	log.Info("[Controller] Calling FolderController.FindFolderContentList")
	rows, _, err := ctl.folderService.FindFolderContentList(nil)
	if err != nil {
		log.Errorf("[FolderController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *FolderController) QueryFolderDefinition() result.RtnDataList {
	log.Info("[Controller] Calling FolderController.GenerateJson")
	rows, _, err := ctl.folderService.GenerateFolderDefinition()
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
