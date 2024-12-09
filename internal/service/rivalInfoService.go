package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalInfoService struct {
	db *gorm.DB
}

func NewRivalInfoService(db *gorm.DB) *RivalInfoService {
	return &RivalInfoService{
		db: db,
	}
}

func (s *RivalInfoService) FindRivalInfoList() ([]*entity.RivalInfo, int, error) {
	var rivals []*entity.RivalInfo
	s.db.Find(&rivals)
	return rivals, len(rivals), nil
}
