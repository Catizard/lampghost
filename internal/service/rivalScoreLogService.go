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

// Extend function to findRivalScoreLogList
//
// Returns sha256 grouped array
func findRivalScoreLogSha256Map(tx *gorm.DB, rivalID uint) (map[string][]entity.RivalScoreLog, error) {
	scoreLogs, _, err := findRivalScoreLogList(tx, rivalID)
	if err != nil {
		return nil, err
	}
	sha256ScoreLogsMap := make(map[string][]entity.RivalScoreLog)
	for _, scoreLog := range scoreLogs {
		if _, ok := sha256ScoreLogsMap[scoreLog.Sha256]; !ok {
			sha256ScoreLogsMap[scoreLog.Sha256] = make([]entity.RivalScoreLog, 0)
		}
		sha256ScoreLogsMap[scoreLog.Sha256] = append(sha256ScoreLogsMap[scoreLog.Sha256], *scoreLog)
	}
	return sha256ScoreLogsMap, nil
}
