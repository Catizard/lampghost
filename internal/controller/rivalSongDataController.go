package controller

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
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

func (ctl *RivalSongDataController) QuerySongDataPageList(param *vo.RivalSongDataVo) result.RtnPage {
	log.Info("[Controller] calling RivalSongDataController.QuerySongDataPageList")
	if param == nil {
		return result.NewErrorPage(fmt.Errorf("QuerySongDataPageList: param should not be empty"))
	}
	rows, _, err := ctl.rivalSongDataService.QuerySongDataPageList(param)
	if err != nil {
		log.Errorf("[RivalSongDataController] returning err: %v", err)
		return result.NewErrorPage(err)
	}
	return result.NewRtnPage(*param.Pagination, rows)
}

func (ctl *RivalSongDataController) ReloadRivalSongData() result.RtnMessage {
	log.Info("[Controller] calling RivalSongDataController.ReloadRivalSongData")
	if err := ctl.rivalSongDataService.ReloadRivalSongData(); err != nil {
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (*RivalSongDataController) GENERATE_RIVAL_SONG_DATA_DTO() *dto.RivalSongDataDto { return nil }
