package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"gorm.io/gorm"
)

type RivalScoreDataLogService struct {
	db *gorm.DB
}

func NewRivalScoreDataLogService(db *gorm.DB) *RivalScoreDataLogService {
	return &RivalScoreDataLogService{
		db: db,
	}
}

func findLastRivalScoreDataLog(tx *gorm.DB, filter *vo.RivalScoreDataLogVo) (*entity.RivalScoreDataLog, error) {
	ret := entity.RivalScoreDataLog{}
	err := tx.Model(&ret).
		Scopes(scopeRivalScoreDataLogFilter(filter)).
		Order("record_time desc").
		Limit(1).
		Find(&ret).
		Error
	return &ret, err
}

// Specialized scope for vo.RivalScoreDataLogVo
func scopeRivalScoreDataLogFilter(filter *vo.RivalScoreDataLogVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter != nil {
			return db
		}
		moved := db.Where(filter.Entity())
		return moved
	}
}
