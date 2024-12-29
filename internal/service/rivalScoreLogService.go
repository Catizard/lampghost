package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalScoreLogService struct {
	db *gorm.DB
}

func NewRivalScoreLogService(db *gorm.DB) *RivalScoreLogService {
	return &RivalScoreLogService{
		db: db,
	}
}

func findRivalScoreLogList(tx *gorm.DB, rivalID uint) ([]*entity.RivalScoreLog, int, error) {
	var out []*entity.RivalScoreLog
	if err := tx.Where("rival_id = ?", rivalID).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, len(out), nil
}
