package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
)

type CustomDiffTableController struct {
	customDiffTableService *service.CustomDiffTableService
}

func NewCustomDiffTableController(customDiffTableService *service.CustomDiffTableService) *CustomDiffTableController {
	return &CustomDiffTableController{
		customDiffTableService: customDiffTableService,
	}
}

func (ctl *CustomDiffTableController) FindCustomDiffTableList(filter *vo.CustomDiffTableVo) result.RtnDataList {
	log.Info("[Controller] Calling CustomDiffTableController.FindCustomDiffTableList")
	rows, _, err := ctl.customDiffTableService.FindCustomDiffTableList(filter)
	if err != nil {
		log.Errorf("[CustomDiffTableController] returning err: %s", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *CustomDiffTableController) FindCustomDiffTableByID(ID uint) result.RtnData {
	log.Info("[Controller] Calling CustomDiffTableController.FindCustomDiffTableByID")
	data, err := ctl.customDiffTableService.FindCustomDiffTableByID(ID)
	if err != nil {
		log.Errorf("[CustomDiffTableController] returning err: %s", err)
		return result.NewErrorData(err)
	}
	return result.NewRtnData(data)
}

func (ctl *CustomDiffTableController) QueryCustomDiffTablePageList(filter *vo.CustomDiffTableVo) result.RtnPage {
	log.Info("[Controller] Calling CustomDiffTableController.QueryCustomDiffTablePageList")
	rows, _, err := ctl.customDiffTableService.FindCustomDiffTableList(filter)
	if err != nil {
		log.Errorf("[CustomDiffTableController] returning err: %s", err)
		return result.NewErrorPage(err)
	}
	return result.NewRtnPage(*filter.Pagination, rows)
}

func (ctl *CustomDiffTableController) UpdateCustomDiffTable(param *vo.CustomDiffTableVo) result.RtnMessage {
	log.Info("[Controller] Calling CustomDiffTableController.UpdateCustomDiffTable")
	if err := ctl.customDiffTableService.UpdateCustomDiffTable(param); err != nil {
		log.Errorf("[CustomDiffTableController] returning err: %s", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *CustomDiffTableController) GENERATOR_CUSTOM_DIFF_TABLE_DTO() *dto.CustomDiffTableDto {
	return nil
}
