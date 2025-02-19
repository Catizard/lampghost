package service

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"gorm.io/gorm"
)

type CourseInfoService struct {
	db *gorm.DB
}

func NewCourseInfoSerivce(db *gorm.DB) *CourseInfoService {
	return &CourseInfoService{
		db: db,
	}
}

func (s *CourseInfoService) FindCourseInfoList(filter *vo.CourseInfoVo) ([]*dto.CourseInfoDto, int, error) {
	return findCourseInfoList(s.db, filter)
}

func findCourseInfoList(tx *gorm.DB, filter *vo.CourseInfoVo) ([]*dto.CourseInfoDto, int, error) {
	if filter == nil {
		var ret []*dto.CourseInfoDto
		if err := tx.Model(&entity.CourseInfo{}).Find(&ret).Error; err != nil {
			return nil, 0, err
		}
		return ret, len(ret), nil
	}

	var ret []*dto.CourseInfoDto
	if err := tx.Model(&entity.CourseInfo{}).Where(filter.Entity()).Find(&ret).Error; err != nil {
		return nil, 0, err
	}
	return ret, len(ret), nil
}
