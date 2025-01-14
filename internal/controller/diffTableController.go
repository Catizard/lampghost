package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/charmbracelet/log"
)

type DiffTableController struct {
	diffTableService *service.DiffTableService
}

func NewDiffTableController(diffTableService *service.DiffTableService) *DiffTableController {
	return &DiffTableController{
		diffTableService: diffTableService,
	}
}

func (ctl *DiffTableController) AddDiffTableHeader(url string) result.RtnMessage {
	log.Info("[Controller] calling DiffTableController.AddDiffTableHeader")
	_, err := ctl.diffTableService.AddDiffTableHeader(url)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *DiffTableController) DelDiffTableHeader(ID uint) result.RtnMessage {
	log.Info("[Controller] calling DiffTableController.DelDiffTableHeader")
	err := ctl.diffTableService.DelDiffTableHeader(ID)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *DiffTableController) FindDiffTableHeaderList() result.RtnDataList {
	log.Info("[Controller] calling DiffTableController.FindDiffTableHeaderList")
	rows, _, err := ctl.diffTableService.FindDiffTableHeaderList()
	if err != nil {
		log.Errorf("[DiffTableController] returning error: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *DiffTableController) FindDiffTableHeaderListWithRival(rivalID uint) result.RtnDataList {
	log.Info("[Controller] calling DiffTableController.FindDiffTableHeaderlistWithRival")
	rows, _, err := ctl.diffTableService.FindDiffTableHeaderListWithRival(rivalID)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *DiffTableController) FindDiffTableHeaderTree() result.RtnDataList {
	log.Info("[Controller] calling DiffTableController.FindDiffTableHeaderTree")
	rows, _, err := ctl.diffTableService.FindDiffTableHeaderTree()
	if err != nil {
		log.Errorf("[DiffTableController] returning error: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *DiffTableController) QueryDiffTableInfoById(ID uint) result.RtnData {
	log.Info("[Controller] calling DiffTableController.QueryDiffTableInfoById")
	data, err := ctl.diffTableService.QueryDiffTableInfoByID(ID)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorData(err)
	}
	return result.NewRtnData(data)
}

func (ctl *DiffTableController) QueryDiffTableInfoWithRival(ID uint, rivalID uint) result.RtnData {
	log.Info("[Controller] calling DiffTableController.QueryDiffTableInfoWithRival")
	data, err := ctl.diffTableService.QueryDiffTableInfoByIDWithRival(ID, rivalID)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorData(err)
	}
	return result.NewRtnData(data)
}

func (ctl *DiffTableController) QueryDiffTableDataWithRival(headerID uint, level string, rivalID uint) result.RtnDataList {
	log.Info("[Controller] calling DiffTableController.QueryDiffTableDataWithRival")
	rows, _, err := ctl.diffTableService.QueryDiffTableDataWithRival(headerID, level, rivalID)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *DiffTableController) QueryLevelLayeredDiffTableInfoById(ID uint) result.RtnData {
	log.Info("[Controller] calling QueryLevelLayeredDiffTableInfoById")
	data, err := ctl.diffTableService.QueryLevelLayeredDiffTableInfoById(ID)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorData(err)
	}
	return result.NewRtnData(data)
}

func (ctl *DiffTableController) GENERATOR_TABLE_HEADER() *entity.DiffTableHeader { return nil }

func (ctl *DiffTableController) GENERATOR_TABLE_HEADER_DTO() *dto.DiffTableHeaderDto { return nil }

func (ctl *DiffTableController) GENERATOR_TABLE_DATA() *dto.DiffTableDataDto { return nil }
