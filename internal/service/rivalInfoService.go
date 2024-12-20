package service

import (
	"fmt"
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"gorm.io/driver/sqlite"
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
	if err := s.db.Find(&rivals).Error; err != nil {
		return nil, 0, err
	}
	return rivals, len(rivals), nil
}

func (s *RivalInfoService) FindRivalInfoByID(rivalID uint) (*entity.RivalInfo, error) {
	out := entity.RivalInfo{}
	if err := s.db.First(&out, rivalID).Error; err != nil {
		log.Debugf("[RivalInfoService] FindRivalInfoByID with ID=%d failed: %v\n", rivalID, err)
		return nil, err
	}
	return &out, nil
}

// TODO: 移除掉该接口同步SongData数据的能力
func (s *RivalInfoService) SyncRivalScoreLog(rivalInfo *entity.RivalInfo) error {
	log.Debug("[Service] calling RivalInfoService.SyncRivalScoreLog")
	if rivalInfo.ScoreLogPath == nil || *rivalInfo.ScoreLogPath == "" {
		return fmt.Errorf("Cannot sync rival %s's score log: score log file path is empty!", rivalInfo.Name)
	}
	log.Debugf("[RivalInfoService] Trying to read score log from %s", *rivalInfo.ScoreLogPath)
	scoreLogDB, err := gorm.Open(sqlite.Open(*rivalInfo.ScoreLogPath))
	if err != nil {
		return err
	}
	scoreLogService := NewScoreLogService(scoreLogDB)
	rawScoreLog, n, err := scoreLogService.FindScoreLogList()
	if err != nil {
		return err
	}
	if n == 0 {
		return nil
	}
	log.Debugf("[RivalInfoService] Read %d logs from %s", n, *rivalInfo.ScoreLogPath)
	log.Debugf("[RivalInfoService] Trying to read song data from %s", *rivalInfo.SongDataPath)
	songDataDB, err := gorm.Open(sqlite.Open(*rivalInfo.SongDataPath))
	if err != nil {
		return err
	}
	songDataService := NewSongDataService(songDataDB)
	rawSongData, n, err := songDataService.FindSongDataList()
	if err != nil {
		return err
	}
	log.Infof("[RivalInfoService] Read %d song data from %s", n, *rivalInfo.SongDataPath)
	if n == 0 {
		return nil
	}

	rivalScoreLog := make([]entity.RivalScoreLog, len(rawScoreLog))
	for i, rawLog := range rawScoreLog {
		rivalLog := entity.FromRawScoreLogToRivalScoreLog(rawLog)
		rivalLog.RivalId = rivalInfo.ID
		rivalScoreLog[i] = rivalLog
	}

	rivalSongData := make([]entity.RivalSongData, len(rawSongData))
	for i, rawData := range rawSongData {
		rivalData := entity.FromRawSongDataToRivalSongData(rawData)
		rivalData.RivalId = rivalInfo.ID
		rivalSongData[i] = rivalData
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("rival_id = ?", rivalInfo.ID).Delete(&entity.RivalScoreLog{}).Error; err != nil {
			return err
		}

		if err := tx.CreateInBatches(&rivalScoreLog, 100).Error; err != nil {
			return err
		}

		if err := tx.CreateInBatches(&rivalSongData, 100).Error; err != nil {
			return err
		}

		if err := tx.Model(rivalInfo).Updates(entity.RivalInfo{
			PlayCount: len(rivalScoreLog),
		}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	log.Debugf("[RivalInfoService] Sync rival %s's %d scorelogs successfully!", rivalInfo.Name, len(rawScoreLog))
	return nil
}

func (s *RivalInfoService) DelRivalInfo(ID uint) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		var candidate entity.RivalInfo
		if err := tx.First(&candidate, ID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&entity.RivalInfo{}, candidate.ID).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *RivalInfoService) QueryUserPlayCountInYear(ID uint, yearNum int) ([]int, error) {
	var playData []entity.RivalScoreLog
	if err := s.db.Where(&entity.RivalScoreLog{RivalId: ID}).Find(&playData).Error; err != nil {
		return nil, err
	}
	ret := make([]int, 12)
	for i := range ret {
		ret[i] = 0
	}
	for _, playLog := range playData {
		date := time.Unix(playLog.TimeStamp, 0)
		ret[date.Month()-1]++
	}
	return ret, nil
}

func (s *RivalInfoService) SyncRivalScoreLogByID(rivalID uint) error {
	if rivalInfo, err := s.FindRivalInfoByID(rivalID); err != nil {
		return err
	} else {
		return s.SyncRivalScoreLog(rivalInfo)
	}
}

func (s *RivalInfoService) QueryMainUser() (*entity.RivalInfo, error) {
	var out entity.RivalInfo
	if err := s.db.Where(&entity.RivalInfo{MainUser: true}).First(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
