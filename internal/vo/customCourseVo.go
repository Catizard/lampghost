package vo

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type CustomCourseVo struct {
	gorm.Model
	Name          string
	CustomTableID uint
	Sha256s       string
	Md5s          string
	Constraints   string
	OrderNumber   int

	// This two fields cannot be translated into entity directly
	SplitSha256s string
	SplitMd5s    string
}

func (customCourse *CustomCourseVo) Entity() *entity.CustomCourse {
	return &entity.CustomCourse{
		Model: gorm.Model{
			ID:        customCourse.ID,
			CreatedAt: customCourse.CreatedAt,
			UpdatedAt: customCourse.UpdatedAt,
			DeletedAt: customCourse.DeletedAt,
		},
		Name:          customCourse.Name,
		CustomTableID: customCourse.CustomTableID,
		Sha256s:       customCourse.Sha256s,
		Md5s:          customCourse.Md5s,
		Constraints:   customCourse.Constraints,
		OrderNumber:   customCourse.OrderNumber,
	}
}
