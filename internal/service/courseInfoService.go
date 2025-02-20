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

// This function returns CourseInfoDto directly due to some historical problem
//
// scorelog.db only records the sha256 while courses doesn't have, what it has is "md5"
// therefore we always need to link the md5 with sha256, so the code is written in basic find method
//
// TODO: This function forces using the main user's songdata.db to build the cache, should be implemented with a config
func findCourseInfoList(tx *gorm.DB, filter *vo.CourseInfoVo) ([]*dto.CourseInfoDto, int, error) {
	partial := tx.Model(&entity.CourseInfo{})
	if filter != nil {
		partial = partial.Where(filter.Entity())
	}
	var raw []*entity.CourseInfo
	if err := partial.Find(&raw).Error; err != nil {
		return nil, 0, err
	}
	if len(raw) == 0 {
		return make([]*dto.CourseInfoDto, 0), 0, nil
	}
	cache, err := queryDefaultSongHashCache(tx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*dto.CourseInfoDto, len(raw))
	for i := range raw {
		out[i] = dto.NewCourseInfoDto(raw[i], cache)
	}
	return out, len(out), nil
}
