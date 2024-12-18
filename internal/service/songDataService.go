package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type SongDataService struct {
	db *gorm.DB
}

func NewSongDataService(db *gorm.DB) *SongDataService {
	return &SongDataService{
		db: db,
	}
}

func (s *SongDataService) FindSongDataList() ([]*entity.SongData, int, error) {
	var data []*entity.SongData
	if err := s.db.Find(&data).Error; err != nil {
		return nil, 0, err
	}
	return data, len(data), nil
}
