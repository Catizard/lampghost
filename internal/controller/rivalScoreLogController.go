package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
)

type RivalScoreLogController struct {
	rivalScoreLogService *service.RivalScoreLogService
}

func NewRivalScoreLogController(rivalScoreLogService *service.RivalScoreLogService) *RivalScoreLogController {
	return &RivalScoreLogController{
		rivalScoreLogService: rivalScoreLogService,
	}
}

func (ctl *RivalScoreLogController) QueryRivalScoreLogPageList(param *vo.RivalScoreLogVo) result.RtnDataList {
	log.Info("[Controller] calling RivalScoreLogController.QueryRivalScoreLogPageList")
	rows, _, err := ctl.rivalScoreLogService.QueryRivalScoreLogPageList(param)
	if err != nil {
		log.Error("[RivalScoreLogController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *RivalScoreLogController) GENERATE_RIVAL_SCORE_LOG() *dto.RivalScoreLogDto { return nil }
