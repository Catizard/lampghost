package controller

import (
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
)

type CustomCourseController struct {
	customCourseService *service.CustomCourseService
}

func NewCustomCourseController(customCourseService *service.CustomCourseService) *CustomCourseController {
	return &CustomCourseController{
		customCourseService: customCourseService,
	}
}

func (ctl *CustomCourseController) AddCustomCourse(param *vo.CustomCourseVo) result.RtnMessage {
	log.Info("[Controller] calling CustomCourseController.AddCustomCourse")
	if err := ctl.customCourseService.AddCustomCourse(param); err != nil {
		log.Errorf("[CustomCourseController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *CustomCourseController) FindCustomCourseList(filter *vo.CustomCourseVo) result.RtnDataList {
	log.Info("[Controller] calling CustomCourseController.FindCustomCourseList")
	rows, _, err := ctl.customCourseService.FindCustomCourseList(filter)
	if err != nil {
		log.Errorf("[CustomCourseController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}
