package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/charmbracelet/log"
)

type RivalSongDataController struct {
	rivalSongDataService *service.RivalSongDataService
}

func NewRivalSongDataController(rivalSongDataService *service.RivalSongDataService) *RivalSongDataController {
	return &RivalSongDataController{
		rivalSongDataService: rivalSongDataService,
	}
}

func (ctl *RivalSongDataController) BindRivalSongDataToFolder(rivalSongDataID uint, folderIDs []uint) result.RtnMessage {
	log.Info("[Controller] Calling RivalSongDataController.BindRivalSongDataToFolder")
	err := ctl.rivalSongDataService.BindRivalSongDataToFolder(rivalSongDataID, folderIDs)
	if err != nil {
		log.Errorf("[RivalSongDataController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (*RivalSongDataController) GENERATE_RIVAL_SONG_DATA_DTO() *dto.RivalSongDataDto { return nil }
