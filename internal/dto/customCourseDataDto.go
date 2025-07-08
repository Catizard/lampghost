package dto

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type CustomCourseDataDto struct {
	gorm.Model
	Sha256         string
	Md5            string
	CustomCourseID uint
	OrderNumber    int
}

func NewCustomCourseDataDto(courseData *entity.CustomCourseData) *CustomCourseDataDto {
	return &CustomCourseDataDto{
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
