package controller

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/charmbracelet/log"
)

type RivalInfoController struct {
	rivalInfoService *service.RivalInfoService
}

func NewRivalInfoController(rivalInfoService *service.RivalInfoService) *RivalInfoController {
	return &RivalInfoController{
		rivalInfoService: rivalInfoService,
	}
}

// TODO: 使用统一的返回类，不要把原始类型丢出controller
func (ctl *RivalInfoController) QueryUserInfoByID(rivalID uint) *entity.RivalInfo {
	log.Info("[Controller] calling RivalInfoController.QueryUserInfo")
	out, _ := ctl.rivalInfoService.FindRivalInfoByID(rivalID)
	return out
}

func (ctl *RivalInfoController) SyncRivalScoreLog(rivalID uint) error {
	log.Info("[Controller] calling RivalInfoController.SyncRivalScorelog")
	return ctl.rivalInfoService.SyncRivalScoreLogByID(rivalID)
}
