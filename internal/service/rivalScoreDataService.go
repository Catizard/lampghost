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

func findRivalScoreDataList(tx *gorm.DB, filter *vo.RivalScoreDataVo) (out []*entity.RivalScoreData, n int, err error) {
	err = eris.Wrap(tx.Model(&entity.RivalScoreData{}).Scopes(scopeRivalScoreData(filter)).Find(&out).Error, "query rival_score_data")
	n = len(out)
	return
}

func findRivalMaximumClearScoreDataList(tx *gorm.DB, rivalID uint) ([]*entity.RivalScoreData, int, error) {
	rows, err := tx.Raw(`select * from (
      select *, ROW_NUMBER() OVER w as rn
	    from rival_score_data rscore
      where rival_id = ?
	    WINDOW w AS (PARTITION BY rscore.sha256 ORDER BY clear desc, minbp asc)
  ) rscore where rscore.rn = 1`, rivalID).Rows()
	if err != nil {
		return nil, 0, eris.Wrap(err, "query rival_score_data")
	}
	var ret []*entity.RivalScoreData
	for rows.Next() {
		var data entity.RivalScoreData
		if err := tx.ScanRows(rows, &data); err != nil {
			return nil, 0, eris.Wrap(err, "scan rows")
		}
		ret = append(ret, &data)
	}
	return ret, len(ret), nil
}

func findLastRivalScoreData(tx *gorm.DB, filter *vo.RivalScoreDataVo) (*entity.RivalScoreData, error) {
	ret := entity.RivalScoreData{}
	err := tx.Model(&ret).Scopes(scopeRivalScoreData(filter)).Order("record_time desc").Limit(1).Find(&ret).Error
	return &ret, eris.Wrap(err, "query rival_score_data")
}

func scopeRivalScoreData(filter *vo.RivalScoreDataVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Extra filter here
		return moved
	}
}
