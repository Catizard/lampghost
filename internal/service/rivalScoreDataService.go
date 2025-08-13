package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type RivalScoreDataService struct {
	db *gorm.DB
}

func NewRivalScoreDataService(db *gorm.DB) *RivalScoreDataService {
	return &RivalScoreDataService{
		db: db,
	}
}

// Fully delete all content from rival_score_data and reinsert them
func syncScoreData(tx *gorm.DB, rivalScoreData []*entity.RivalScoreData, rivalID uint) error {
	if err := tx.Unscoped().Where("rival_id = ?", rivalID).Delete(&entity.RivalScoreData{}).Error; err != nil {
		return eris.Wrap(err, "delete rival_score_data")
	}
	return eris.Wrap(tx.CreateInBatches(&rivalScoreData, DEFAULT_BATCH_SIZE).Error, "insert rival_score_data")
}

// Similar to syncScoreData but not delete any old content, only update new logs
func updateScoreData(tx *gorm.DB, rivalScoreData []*entity.RivalScoreData) error {
	for _, data := range rivalScoreData {
		if err := tx.Save(data).Error; err != nil {
			return eris.Wrap(err, "insertOrUpdate rival_score_data")
		}
	}
	return nil
}

func findLastRivalScoreData(tx *gorm.DB, filter *vo.RivalScoreDataVo) (*entity.RivalScoreData, error) {
	ret := entity.RivalScoreData{}
	err := tx.Model(&ret).Order("record_time desc").Limit(1).Find(&ret).Error
	return &ret, eris.Wrap(err, "query rival_score_data")
}
