package controller

import (
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

func (ctl *RivalTagController) QueryRivalTagPageList(filter *vo.RivalTagVo) result.RtnPage {
	log.Info("[Controller] Calling RivalTagController.QueryRivalTagPageList")
	rows, _, err := ctl.rivalTagService.FindRivalTagList(filter)
	if err != nil {
		log.Errorf("[RivalTagController] returning err: %v", err)
		return result.NewErrorPage(err)
	}
	return result.NewRtnPage(*filter.Pagination, rows)
}

func (*RivalTagController) GENERATE_RIVAL_TAG() *entity.RivalTag     { return nil }
func (*RivalTagController) GENERATE_RIVAL_TAG_DTO() *dto.RivalTagDto { return nil }
