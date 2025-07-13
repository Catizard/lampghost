package controller

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
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

func (ctl *CustomCourseController) FindCustomCourseByID(customCourseID uint) result.RtnData {
	log.Info("[Controller] calling CustomCourseController.FindCustomCourseByID")
	data, err := ctl.customCourseService.FindCustomCourseByID(customCourseID)
	if err != nil {
		log.Errorf("[CustomCourseController] returning err: %v", err)
		return result.NewErrorData(err)
	}
	return result.NewRtnData(data)
}

func (ctl *CustomCourseController) QueryCustomCourseSongListWithRival(filter *vo.CustomCourseVo) result.RtnDataList {
	log.Info("[Controller] calling CustomCourseController.QueryCustomCourseSongListWithRival")
	rows, _, err := ctl.customCourseService.QueryCustomCourseSongListWithRival(filter)
	if err != nil {
		log.Errorf("[CustomCourseController] returning err: %v", err)
		return result.NewErrorDataList(err)
	}
	return result.NewRtnDataList(rows)
}

func (ctl *CustomCourseController) AddCustomCourseData(param *entity.CustomCourseData) result.RtnMessage {
	log.Info("[Controller] calling CustomCourseController.AddCustomCourseData")
	if err := ctl.customCourseService.AddCustomCourseData(param); err != nil {
		log.Errorf("[CustomCourseController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *CustomCourseController) BindSongToCustomCourse(sha256, md5 string, customCourseID uint) result.RtnMessage {
	log.Info("[Controller] calling CustomCourseController.BindSongToCustomCourse")
	if err := ctl.customCourseService.BindSongToCustomCourse(sha256, md5, customCourseID); err != nil {
		log.Errorf("[CustomCourseController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *CustomCourseController) UpdateCustomCourse(param *vo.CustomCourseVo) result.RtnMessage {
	log.Info("[Controller] calling CustomCourseController.UpdateCustomCourse")
	if err := ctl.customCourseService.UpdateCustomCourse(param); err != nil {
		log.Errorf("[CustomCourseController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *CustomCourseController) UpdateCustomCourseOrder(customCourseIDs []uint) result.RtnMessage {
	log.Info("[Controller] calling CustomCourseController.UpdateCustomCourseOrder")
	if err := ctl.customCourseService.UpdateCustomCourseOrder(customCourseIDs); err != nil {
		log.Errorf("[CustomCourseController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *CustomCourseController) UpdateCustomCourseDataOrder(customCourseDataIDs []uint) result.RtnMessage {
	log.Info("[Controller] calling CustomCourseController.UpdateCustomCourseDataOrder")
	if err := ctl.customCourseService.UpdateCustomCourseDataOrder(customCourseDataIDs); err != nil {
		log.Errorf("[CustomCourseController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *CustomCourseController) DeleteCustomCourse(courseID uint) result.RtnMessage {
	log.Info("[Controller] calling CustomCourseController.DeleteCustomCourse")
	if err := ctl.customCourseService.DeleteCustomCourse(courseID); err != nil {
		log.Errorf("[CustomCourseController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *CustomCourseController) DeleteCustomCourseData(courseDataID uint) result.RtnMessage {
	log.Info("[Controller] calling CustomCourseController.UpdateCustomCourseDataOrder")
	if err := ctl.customCourseService.DeleteCustomCourseData(courseDataID); err != nil {
		log.Errorf("[CustomCourseController] returning err: %v", err)
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func (ctl *CustomCourseController) GENERATOR_CUSTOM_COUSE() *entity.CustomCourse { return nil }
func (ctl *CustomCourseController) GENERATOR_CUSTOM_COUSE_DTO() *dto.CustomCourseDto          { return nil }
func (ctl *CustomCourseController) GENERATOR_CUSTOM_COUSE_DATA() *dto.CustomCourseDataDto { return nil }
