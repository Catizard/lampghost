package vo

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type CustomCourseDataVo struct {
	gorm.Model
	Sha256         string
	Md5            string
	CustomCourseID uint
	OrderNumber    int

	CustomCourseIDs []uint
}

func (courseData *CustomCourseDataVo) Entity() *entity.CustomCourseData {
	return &entity.CustomCourseData{
		Model: gorm.Model{
			ID:        courseData.ID,
			CreatedAt: courseData.CreatedAt,
			UpdatedAt: courseData.UpdatedAt,
			DeletedAt: courseData.DeletedAt,
		},
		Sha256:         courseData.Sha256,
		Md5:            courseData.Md5,
		CustomCourseID: courseData.CustomCourseID,
		OrderNumber:    courseData.OrderNumber,
	}
}
