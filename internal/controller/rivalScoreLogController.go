package controller

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/rotisserie/eris"
)

type RivalScoreLogController struct {
	rivalScoreLogService *service.RivalScoreLogService
}

func NewRivalScoreLogController(rivalScoreLogService *service.RivalScoreLogService) *RivalScoreLogController {
	return &RivalScoreLogController{
		rivalScoreLogService: rivalScoreLogService,
	}
}

func (ctl *RivalScoreLogController) QueryRivalScoreLogPageList(param *vo.RivalScoreLogVo) result.RtnPage {
	log.Info("[Controller] calling RivalScoreLogController.QueryRivalScoreLogPageList")
	if param == nil {
		return result.NewErrorPage(fmt.Errorf("QueryRivalScoreLogPageList: param should not be empty"))
	}
	rows, _, err := ctl.rivalScoreLogService.QueryRivalScoreLogPageList(param)
	if err != nil {
		log.Errorf("[RivalScoreLogController] returning err: %v", err)
		return result.NewErrorPage(err)
	}
	return result.NewRtnPage(*param.Pagination, rows)
}

func (ctl *RivalScoreLogController) QueryPrevDayScoreLogList(param *vo.RivalScoreLogVo) result.RtnDataList {
	log.Info("[Controller] calling RivalScoreLogController.QueryPrevDayScoreLogList")
	if param == nil {
		return result.NewErrorDataList(eris.Errorf("QueryPrevDayScoreLogList: param should not be empty"))
	}
	param.ConvTimestamp()
	rows, _, err := ctl.rivalScoreLogService.QueryPrevDayScoreLogList(param)
	if err != nil {
		log.Errorf("[RivalScoreLogController] returning err: %v", err)
		return result.NewErrorDataList(eris.Cause(err))
	}
	return result.NewRtnDataList(rows)
}

func (ctl *RivalScoreLogController) GENERATE_RIVAL_SCORE_LOG() *dto.RivalScoreLogDto { return nil }
