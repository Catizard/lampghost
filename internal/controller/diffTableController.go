package controller

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
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

func (ctl *DiffTableController) AddDiffTableHeader(url string) *entity.DiffTableHeader {
	log.Info("[Controller] calling DiffTableController.AddDiffTableHeader")
	ret, err := ctl.diffTableService.AddDiffTableHeader(url)
	if err != nil {
		log.Errorf("[DiffTableController] returning err: %v", err)
	}
	return ret
}

func (ctl *DiffTableController) FindDiffTableHeader() []*entity.DiffTableHeader {
	log.Info("[Controller] calling DiffTableController.FindDiffTableHeader")
	ret, _, err := ctl.diffTableService.FindDiffTableHeader()
	if err != nil {
		log.Errorf("[DiffTableController] returning error: %v", err)
	}
	return ret
}
