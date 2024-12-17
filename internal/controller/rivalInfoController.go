package controller

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/result"
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

func (ctl *RivalInfoController) QueryUserInfoByID(rivalID uint) result.RtnData {
	log.Info("[Controller] calling RivalInfoController.QueryUserInfo")
	data, err := ctl.rivalInfoService.FindRivalInfoByID(rivalID)
	if err != nil {
		return result.NewErrorData(err)
	}
	return result.NewRtnData(data)
}

func (ctl *RivalInfoController) SyncRivalScoreLog(rivalID uint) error {
	log.Info("[Controller] calling RivalInfoController.SyncRivalScorelog")
	return ctl.rivalInfoService.SyncRivalScoreLogByID(rivalID)
}

// TODO: 目前wails不支持泛型代码生成，这个方法用于让wails知道需要生成entity下的数据
func (ctl *RivalInfoController) GENERATOR_RIVAL_INFO() *entity.RivalInfo { return nil }
