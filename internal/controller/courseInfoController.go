package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/charmbracelet/log"
)

type CourseInfoController struct {
	courseInfoService *service.CourseInfoService
}

func NewCourseInfoController(courseInfoService *service.CourseInfoService) *CourseInfoController {
	return &CourseInfoController{
		courseInfoService: courseInfoService,
	}
}

func (ctl *CourseInfoController) FindCourseInfoList() result.RtnDataList {
	log.Info("[Controller] calling CourseInfoController.FindCourseInfoList")
	rows, _, err := ctl.courseInfoService.FindCourseInfoList(nil)
	if err != nil {
		log.Errorf("[CourseInfoController] returning error: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (*CourseInfoController) GENERATOR_COURSE_INFO() *entity.CourseInfo     { return nil }
func (*CourseInfoController) GENERATOR_COURSE_INFO_DTO() *dto.CourseInfoDto { return nil }
