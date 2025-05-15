package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/rotisserie/eris"
)

type RivalScoreDataLogController struct {
	rivalScoreDataLogService *service.RivalScoreDataLogService
}

func NewRivalScoreDataLogController(rivalScoreDataLogService *service.RivalScoreDataLogService) *RivalScoreDataLogController {
	return &RivalScoreDataLogController{
		rivalScoreDataLogService: rivalScoreDataLogService,
	}
}

func (ctl *RivalScoreDataLogController) QueryUserKeyCountInYear(param *vo.RivalScoreDataLogVo) result.RtnDataList {
	log.Info("[Controller] calling RivalScoreDataLogController.QueryUserKeyCountInYear")
	if param == nil {
		return result.NewErrorDataList(eris.Errorf("QueryPrevDayScoreLogList: param should not be empty"))
	}
	rows, _, err := ctl.rivalScoreDataLogService.QueryUserKeyCountInYear(param)
	if err != nil {
		log.Errorf("[RivalScoreDataLogController] returning err: %v", err)
		return result.NewErrorDataList(eris.Cause(err))
	}
	return result.NewRtnDataList(rows)

}

func (ctl *RivalScoreDataLogController) GENERATOR_KEY_COUNT_DTO() *dto.KeyCountDto { return nil }
