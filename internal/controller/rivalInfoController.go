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

func (ctl *RivalInfoController) InitializeMainUser(rivalInfo *entity.RivalInfo) result.RtnMessage {
	log.Info("[Controller] calling RivalInfoController.InitializeMainUser")
	err := ctl.rivalInfoService.InitializeMainUser(rivalInfo)
	if err != nil {
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *RivalInfoController) QueryMainUser() result.RtnData {
	log.Info("[Controller] calling RivalInfoController.QueryMainUser")
	rivalInfo, err := ctl.rivalInfoService.QueryMainUser()
	if err != nil {
		return result.NewErrorData(err)
	}
	return result.NewRtnData(rivalInfo)
}

func (ctl *RivalInfoController) QueryUserInfoByID(rivalID uint) result.RtnData {
	log.Info("[Controller] calling RivalInfoController.QueryUserInfo")
	data, err := ctl.rivalInfoService.FindRivalInfoByID(rivalID)
	if err != nil {
		return result.NewErrorData(err)
	}
	return result.NewRtnData(data)
}

func (ctl *RivalInfoController) QueryUserPlayCountInYear(ID uint, yearNum int) result.RtnDataList {
	log.Info("[Controller] calling RivalInfoController.QueryUserPlayCountInYear")
	pc, err := ctl.rivalInfoService.QueryUserPlayCountInYear(ID, yearNum)
	if err != nil {
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(pc)
}

func (ctl *RivalInfoController) QueryUserInfoWithLevelLayeredDiffTableLampStatus(rivalID uint, headerID uint) result.RtnData {
	log.Info("[Controller] calling RivalInfoController.QueryDiffTableLampStatus")
	data, err := ctl.rivalInfoService.QueryUserInfoWithLevelLayeredDiffTableLampStatus(rivalID, headerID)
	if err != nil {
		return result.NewErrorData(err)
	}
	return result.NewRtnData(data)
}

func (ctl *RivalInfoController) SyncRivalScoreLog(rivalID uint) result.RtnMessage {
	log.Info("[Controller] calling RivalInfoController.SyncRivalScorelog")
	if err := ctl.rivalInfoService.SyncRivalScoreLogByID(rivalID); err != nil {
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *RivalInfoController) FindRivalInfoList() result.RtnDataList {
	log.Info("[Controller] calling RivalInfoController.FindRivalInfoList")
	rows, _, err := ctl.rivalInfoService.FindRivalInfoList()
	if err != nil {
		log.Errorf("[RivalInfoController] returning error: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *RivalInfoController) DelRivalInfo(ID uint) result.RtnMessage {
	log.Info("[Controller] calling RivalInfoController.DelRivalInfo")
	err := ctl.rivalInfoService.DelRivalInfo(ID)
	if err != nil {
		log.Errorf("[RivalInfoController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

// TODO: 目前wails不支持泛型代码生成，这个方法用于让wails知道需要生成entity下的数据
func (ctl *RivalInfoController) GENERATOR_RIVAL_INFO() *entity.RivalInfo { return nil }
