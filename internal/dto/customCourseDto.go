package dto

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type CustomCourseDto struct {
	gorm.Model
	Name          string
	CustomTableID uint
	Sha256s       string
	Md5s          string
	Constraints   string
	OrderNumber   int
}

func (customCourse *CustomCourseDto) Entity() *entity.CustomCourse {
	return &entity.CustomCourse{
		Model: gorm.Model{
			ID:        customCourse.ID,
			CreatedAt: customCourse.CreatedAt,
			UpdatedAt: customCourse.UpdatedAt,
			DeletedAt: customCourse.DeletedAt,
		},
		Name:          customCourse.Name,
		CustomTableID: customCourse.CustomTableID,
		Constraints:   customCourse.Constraints,
		OrderNumber:   customCourse.OrderNumber,
	}
}
