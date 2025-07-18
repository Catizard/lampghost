package controller

import (
	"context"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
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

func (ctl *RivalInfoController) InjectContext(ctx context.Context) {
	ctl.rivalInfoService.InjectContext(ctx)
}

func (ctl *RivalInfoController) InitializeMainUser(rivalInfo *vo.InitializeRivalInfoVo) result.RtnMessage {
	log.Info("[Controller] calling RivalInfoController.InitializeMainUser")
	err := ctl.rivalInfoService.InitializeMainUser(rivalInfo)
	if err != nil {
		log.Errorf("[RivalController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *RivalInfoController) ChooseBeatorajaDirectory() result.RtnData {
	log.Info("[Controller] calling RivalInfoController.ChooseBeatorajaDirectory")
	meta, err := ctl.rivalInfoService.ChooseBeatorajaDirectory()
	if err != nil {
		log.Errorf("[RivalInfoController] returning err: %v", err)
		return result.NewErrorData(err)
	}
	return result.NewRtnData(meta)
}

// TODO: Seperate songdata.db file feature doesn't implement yet, this function would
// not save the songdata.db file path parameter
func (ctl *RivalInfoController) AddRivalInfo(rivalInfo *vo.RivalInfoVo) result.RtnMessage {
	log.Info("[Controller] calling RivalInfoController.AddRivalInfo")
	err := ctl.rivalInfoService.AddRivalInfo(rivalInfo)
	if err != nil {
		log.Errorf("[RivalController] returning err: %v", err)
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

func (ctl *RivalInfoController) QueryUserPlayCountInYear(ID uint, yearNum string) result.RtnDataList {
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

func (ctl *RivalInfoController) ReloadRivalData(rivalID uint, fullyReload bool) result.RtnMessage {
	log.Info("[Controller] calling RivalInfoController.ReloadRivalData")
	if err := ctl.rivalInfoService.ReloadRivalData(rivalID, fullyReload); err != nil {
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *RivalInfoController) FindRivalInfoList() result.RtnDataList {
	log.Info("[Controller] calling RivalInfoController.FindRivalInfoList")
	rows, _, err := ctl.rivalInfoService.FindRivalInfoList(nil)
	if err != nil {
		log.Errorf("[RivalInfoController] returning error: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *RivalInfoController) QueryRivalInfoPageList(filter *vo.RivalInfoVo) result.RtnPage {
	log.Info("[Controller] calling RivalInfoController.QueryRivalInfoPageList")
	rows, _, err := ctl.rivalInfoService.FindRivalInfoList(filter)
	if err != nil {
		log.Errorf("[RivalInfoController] returning error: %v", err)
		return result.NewErrorPage(err)
	}
	return result.NewRtnPage(*filter.Pagination, rows)
}

func (ctl *RivalInfoController) QueryRivalPlayedYears(rivalID uint) result.RtnDataList {
	log.Info("[Controller] calling RivalInfoController.QueryRivalPlayedYears")
	rows, _, err := ctl.rivalInfoService.QueryRivalPlayedYears(rivalID)
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

func (ctl *RivalInfoController) UpdateRivalInfo(updateParam *vo.RivalInfoVo) result.RtnMessage {
	log.Info("[Controller] calling RivalInfoController.UpdateRivalInfo")
	err := ctl.rivalInfoService.UpdateRivalInfo(updateParam)
	if err != nil {
		log.Errorf("[RivalInfoController] returning error: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *RivalInfoController) UpdateRivalReverseImportInfo(updateParam *vo.RivalInfoVo) result.RtnMessage {
	log.Info("[Controller] calling RivalInfoController.UpdateRivalReverseImportInfo")
	if err := ctl.rivalInfoService.UpdateRivalReverseImportInfo(updateParam); err != nil {
		log.Errorf("[RivalInfoController] returning error: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *RivalInfoController) GENERATOR_RIVAL_INFO() *entity.RivalInfo     { return nil }
func (ctl *RivalInfoController) GENERATOR_RIVAL_INFO_DTO() *dto.RivalInfoDto { return nil }
func (ctl *RivalInfoController) GENERATOR_BEATORAJA_DIRECTORY_META() *dto.BeatorajaDirectoryMeta {
	return nil
}
