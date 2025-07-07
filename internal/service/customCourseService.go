package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type CustomCourseService struct {
	db *gorm.DB
}

func NewCustomCourseService(db *gorm.DB) *CustomCourseService {
	return &CustomCourseService{
		db: db,
	}
}

func (s *CustomCourseService) AddCustomCourse(param *vo.CustomCourseVo) error {
	if param == nil {
		return eris.Errorf("AddCustomCourse: param cannot be nil")
	}
	if param.Name == "" {
		return eris.Errorf("AddCustomCourse: name cannot be empty")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		return addCustomCourse(tx, param)
	})
}

func (s *CustomCourseService) FindCustomCourseList(filter *vo.CustomCourseVo) (out []*entity.CustomCourse, n int, err error) {
	err = s.db.Transaction(func(tx *gorm.DB) error {
		out, n, err = FindCustomCourseList(tx, filter)
		return err
	})
	return
}

func FindCustomCourseList(tx *gorm.DB, filter *vo.CustomCourseVo) (out []*entity.CustomCourse, n int, err error) {
	err = tx.Model(&entity.CustomCourse{}).Scopes(scopeCustomCourseFilter(filter)).Order("order_number").Find(&out).Error
	return
}

func addCustomCourse(tx *gorm.DB, param *vo.CustomCourseVo) error {
	return tx.Create(param.Entity()).Error
}

func scopeCustomCourseFilter(filter *vo.CustomCourseVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Extra filters here
		return moved
	}
}
