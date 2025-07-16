package controller

import (
	"time"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
)

type RivalTagController struct {
	rivalTagService *service.RivalTagService
}

func NewRivalTagController(rivalTagService *service.RivalTagService) *RivalTagController {
	return &RivalTagController{
		rivalTagService: rivalTagService,
	}
}

func (ctl *RivalTagController) FindRivalTagList(filter *vo.RivalTagVo) result.RtnDataList {
	log.Info("[Controller] Calling RivalTagController.FindRivalTagList")
	rows, _, err := ctl.rivalTagService.FindRivalTagList(filter)
	if err != nil {
		log.Errorf("[RivalTagController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *RivalTagController) FindRivalTagByID(rivalTagID uint) result.RtnData {
	log.Info("[Controller] calling RivalTagController.FindRivalTagByID")
	if data, err := ctl.rivalTagService.FindRivalTagByID(rivalTagID); err != nil {
		log.Errorf("[RivalTagController] returning err: %v", err)
		return result.NewErrorData(err)
	} else {
		return result.NewRtnData(data)
	}
}

func (ctl *RivalTagController) QueryRivalTagPageList(filter *vo.RivalTagVo) result.RtnPage {
	log.Info("[Controller] Calling RivalTagController.QueryRivalTagPageList")
	rows, _, err := ctl.rivalTagService.FindRivalTagList(filter)
	if err != nil {
		log.Errorf("[RivalTagController] returning err: %v", err)
		return result.NewErrorPage(err)
	}
	return result.NewRtnPage(*filter.Pagination, rows)
}

func (ctl *RivalTagController) AddRivalTag(rivalTag *vo.RivalTagVo) result.RtnMessage {
	log.Info("[Controller] Calling RivalTagController.AddRivalTag")
	if rivalTag.RecordTimestamp != nil {
		rivalTag.RecordTime = time.Unix((*rivalTag.RecordTimestamp)/1000, 0)
	}
	err := ctl.rivalTagService.AddRivalTag(rivalTag)
	if err != nil {
		log.Errorf("[RivalTagController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *RivalTagController) DeleteRivalTagByID(rivalTagID uint) result.RtnMessage {
	log.Info("[Controller] Calling RivalTagController.DeleteRivalTag")
	err := ctl.rivalTagService.DeleteRivalTagByID(rivalTagID)
	if err != nil {
		log.Errorf("[RivalTagController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *RivalTagController) UpdateRivalTag(param *vo.RivalTagUpdateParam) result.RtnMessage {
	log.Info("[Controller] calling RivalTagController.UpdateRivalTag")
	if err := ctl.rivalTagService.UpdateRivalTag(param); err != nil {
		log.Errorf("[RivalTagController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *RivalTagController) SyncRivalTag(rivalID uint) result.RtnMessage {
	log.Info("[Controller] Calling RivalTagController.SyncRivalTag")
	err := ctl.rivalTagService.SyncRivalTag(rivalID)
	if err != nil {
		log.Error("[RivalTagController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (*RivalTagController) GENERATE_RIVAL_TAG() *entity.RivalTag     { return nil }
func (*RivalTagController) GENERATE_RIVAL_TAG_DTO() *dto.RivalTagDto { return nil }
