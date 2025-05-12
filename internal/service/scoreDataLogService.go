package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type ScoreDataLogService struct {
	db *gorm.DB
}

func NewScoreDataLogService(db *gorm.DB) *ScoreDataLogService {
	return &ScoreDataLogService{
		db: db,
	}
}

func (s *ScoreDataLogService) FindScoreDataLogList(maximumTimestamp *int64) ([]*entity.ScoreDataLog, int, error) {
	var data []*entity.ScoreDataLog
	moved := s.db
	if maximumTimestamp != nil {
		moved = moved.Where("date > ?", maximumTimestamp)
	}
	if err := moved.Find(&data).Error; err != nil {
		return nil, 0, err
	}
	return data, len(data), nil
}
