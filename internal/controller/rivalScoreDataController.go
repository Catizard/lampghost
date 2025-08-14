package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/service"
)

type RivalScoreDataController struct {
	rivalScoreDataService *service.RivalScoreDataService
}

func NewRivalScoreDataController(rivalScoreDataService *service.RivalScoreDataService) *RivalScoreDataController {
	return &RivalScoreDataController{
		rivalScoreDataService: rivalScoreDataService,
	}
}

func (ctl *RivalScoreDataController) GENERATOR_RIVAL_SCORE_DATA_ENTITY() *entity.RivalScoreData {
	return nil
}

func (ctl *RivalScoreDataController) GENERATOR_RIVAL_SCORE_DATA_DTO() *dto.RivalScoreDataDto {
	return nil
}
