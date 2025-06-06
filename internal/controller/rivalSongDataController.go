package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/service"
)

type RivalSongDataController struct {
	rivalSongDataService *service.RivalSongDataService
}

func NewRivalSongDataController(rivalSongDataService *service.RivalSongDataService) *RivalSongDataController {
	return &RivalSongDataController{
		rivalSongDataService: rivalSongDataService,
	}
}

func (*RivalSongDataController) GENERATE_RIVAL_SONG_DATA_DTO() *dto.RivalSongDataDto { return nil }
