package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/rotisserie/eris"
)

type DiffTableController struct {
	diffTableService *service.DiffTableService
}

func NewDiffTableController(diffTableService *service.DiffTableService) *DiffTableController {
	return &DiffTableController{
		diffTableService: diffTableService,
	}
}

func (ctl *DiffTableController) AddDiffTableHeader(param *vo.DiffTableHeaderVo) result.RtnMessage {
	log.Info("[Controller] calling DiffTableController.AddDiffTableHeader")
	_, err := ctl.diffTableService.AddDiffTableHeader(param)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *DiffTableController) AddBatchDiffTableHeader(candidates []*vo.DiffTableHeaderVo) result.RtnDataList {
	log.Info("[Controller] calling DiffTableController.AddBatchDiffTableHeader")
	rows, _, err := ctl.diffTableService.AddBatchDiffTableHeader(candidates)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *DiffTableController) ReloadDiffTableHeader(ID uint) result.RtnMessage {
	log.Info("[Controller] calling DiffTableController.ReloadDiffTableHeader")
	if err := ctl.diffTableService.ReloadDiffTableHeader(ID); err != nil {
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

func (ctl *DiffTableController) UpdateDiffTableHeader(param *vo.DiffTableHeaderVo) result.RtnMessage {
	log.Info("[Controller] calling DiffTableController.UpdateDiffTableHeader")
	err := ctl.diffTableService.UpdateDiffTableHeader(param)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *DiffTableController) FindDiffTableHeaderList() result.RtnDataList {
	log.Info("[Controller] calling DiffTableController.FindDiffTableHeaderList")
	rows, _, err := ctl.diffTableService.FindDiffTableHeaderList(nil)
	if err != nil {
		log.Errorf("[DiffTableController] returning error: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *DiffTableController) FindDiffTableLevelList(ID uint) result.RtnDataList {
	log.Info("[Controller] calling DiffTableController.FindDiffTableLevelList")
	rows, _, err := ctl.diffTableService.FindDiffTableLevelList(ID)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", eris.ToString(err, true))
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *DiffTableController) FindDiffTableHeaderTree(filter *vo.DiffTableHeaderVo) result.RtnDataList {
	log.Info("[Controller] calling DiffTableController.FindDiffTableHeaderTree")
	rows, _, err := ctl.diffTableService.FindDiffTableHeaderTree(filter)
	if err != nil {
		log.Errorf("[DiffTableController] returning error: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *DiffTableController) FindDiffTableHeaderTreeWithRival(filter *vo.DiffTableHeaderVo) result.RtnDataList {
	log.Info("[Controller] calling DiffTableController.FindDiffTableHeaderTreeWithRival")
	rows, _, err := ctl.diffTableService.FindDiffTableHeaderTreeWithRival(filter)
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

func (ctl *DiffTableController) QueryDiffTableDataWithRival(filter *vo.DiffTableHeaderVo) result.RtnPage {
	log.Info("[Controller] calling DiffTableController.QueryDiffTableDataWithRival")
	rows, _, err := ctl.diffTableService.QueryDiffTableDataWithRival(filter)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", eris.ToJSON(err, true))
		return result.NewErrorPage(err)
	}
	return result.NewRtnPage(*filter.Pagination, rows)
}

func (ctl *DiffTableController) UpdateHeaderOrder(headerIDs []uint) result.RtnMessage {
	log.Info("[Controller] Calling DiffTableController.UpdateHeaderOrder")
	if err := ctl.diffTableService.UpdateHeaderOrder(headerIDs); err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *DiffTableController) UpdateHeaderLevelOrders(updateParam *vo.DiffTableHeaderVo) result.RtnMessage {
	log.Info("[Controller] Calling DiffTableController.UpdateHeaderLevelOrders")
	if err := ctl.diffTableService.UpdateHeaderLevelOrders(updateParam); err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *DiffTableController) QueryPredefineTableSchemes() result.RtnDataList {
	log.Info("[Controller] Calling DiffTableController.QueryPredefineTableSchemes")
	schemes := make([]entity.PredefineTableScheme, 0)
	for _, scheme := range entity.PredefineTableSchemes {
		schemes = append(schemes, scheme)
	}
	return result.NewRtnDataList(schemes)
}

func (ctl *DiffTableController) SupplyMissingBMSFromTable(ID uint, levels []string) result.RtnMessage {
	log.Info("[Controller] Calling DiffTableController.SupplyMissingBMSFromTable")
	if err := ctl.diffTableService.SupplyMissingBMSFromTable(ID, levels); err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *DiffTableController) GENERATOR_TABLE_HEADER() *entity.DiffTableHeader { return nil }

func (ctl *DiffTableController) GENERATOR_TABLE_HEADER_DTO() *dto.DiffTableHeaderDto { return nil }

func (ctl *DiffTableController) GENERATOR_TABLE_DATA() *dto.DiffTableDataDto { return nil }

func (ctl *DiffTableController) GENERATOR_TABLE_TAG_DTO() *dto.DiffTableTagDto { return nil }

func (ctl *DiffTableController) GENERATOR_PREDEFINE_TABLE_HEADER() *entity.PredefineTableHeader {
	return nil
}

func (ctl *DiffTableController) GENERATOR_PREDEFINE_TABLE_SCHEME() *entity.PredefineTableScheme {
	return nil
}
