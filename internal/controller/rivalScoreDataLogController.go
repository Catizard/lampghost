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

func (ctl *RivalScoreDataLogController) QueryRivalScoreDataLogPageList(param *vo.RivalScoreDataLogVo) result.RtnPage {
	log.Info("[Controller] calling RivalScoreDataLogController.QueryRivalScoreDataLogPageList")
	if param == nil {
		return result.NewErrorPage(fmt.Errorf("QueryRivalScoreDataLogPageList: param should not be empty"))
	}
	rows, _, err := ctl.rivalScoreDataLogService.QueryRivalScoreDataLogPageList(param)
	if err != nil {
		log.Errorf("[RivalScoreDataLogController] returning err: %v", eris.ToString(err, true))
		return result.NewErrorPage(err)
	}
	return result.NewRtnPage(*param.Pagination, rows)
}

func (ctl *RivalScoreDataLogController) GENERATOR_KEY_COUNT_DTO() *dto.KeyCountDto { return nil }
func (ctl *RivalScoreDataLogController) GENERATOR_RIVAL_SCORE_DATA_LOG_DTO() *dto.RivalScoreDataLogDto {
	return nil
}
