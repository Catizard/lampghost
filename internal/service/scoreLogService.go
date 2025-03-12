package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type ScoreLogService struct {
	db *gorm.DB
}

func NewScoreLogService(db *gorm.DB) *ScoreLogService {
	return &ScoreLogService{
		db: db,
	}
}

func (s *ScoreLogService) FindScoreLogList(maximumTimestamp *int64) ([]*entity.ScoreLog, int, error) {
	var logs []*entity.ScoreLog
	moved := s.db
	if maximumTimestamp != nil {
		moved = moved.Where("date <= ?", maximumTimestamp)
	}
	if err := moved.Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, len(logs), nil
}
