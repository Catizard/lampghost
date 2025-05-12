package controller

import (
	"github.com/Catizard/lampghost_wails/internal/service"
)

type RivalScoreDataLogController struct {
	rivalScoreDataLogService *service.RivalScoreDataLogService
}

func NewRivalScoreDataLogController(rivalScoreDataLogService *service.RivalScoreDataLogService) *RivalScoreDataLogController {
	return &RivalScoreDataLogController{
		rivalScoreDataLogService: rivalScoreDataLogService,
	}
}
