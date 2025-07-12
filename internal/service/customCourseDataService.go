package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"gorm.io/gorm"
)

type CustomCourseDataService struct {
	db *gorm.DB
}

func NewCustomCourseDataService(db *gorm.DB) *CustomCourseService {
	return &CustomCourseService{db: db}
}

func findCustomCourseDataList(tx *gorm.DB, filter *vo.CustomCourseDataVo) (out []*entity.CustomCourseData, n int, err error) {
	err = tx.Model(&entity.CustomCourseData{}).Scopes(scopeCustomCourseDataFilter(filter)).Order("order_number").Find(&out).Error
	n = len(out)
	return
}

func scopeCustomCourseDataFilter(filter *vo.CustomCourseDataVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Extra filters here
		if len(filter.CustomCourseIDs) > 0 {
			moved.Where("custom_course_id in (?)", filter.CustomCourseIDs)
		}
		return moved
	}
}
