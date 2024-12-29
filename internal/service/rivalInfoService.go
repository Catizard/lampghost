package service

import (
	"fmt"
	"time"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RivalInfoService struct {
	db               *gorm.DB
	diffTableService *DiffTableService
	rivalTagService  *RivalTagService
}

func NewRivalInfoService(db *gorm.DB, diffTableService *DiffTableService, rivalTagService *RivalTagService) *RivalInfoService {
	return &RivalInfoService{
		db:               db,
		diffTableService: diffTableService,
		rivalTagService:  rivalTagService,
	}
}

func (s *RivalInfoService) InitializeMainUser(rivalInfo *entity.RivalInfo) error {
	if rivalInfo.SongDataPath == nil || *rivalInfo.SongDataPath == "" {
		return fmt.Errorf("songdata.db path cannot be empty")
	}
	if rivalInfo.ScoreLogPath == nil || *rivalInfo.ScoreLogPath == "" {
		return fmt.Errorf("scorelog.db path cannot be empty")
	}
	var cnt int64
	if err := s.db.Model(&entity.RivalInfo{}).Where("main_user = 1").Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return fmt.Errorf("cannot have two main user, what are you doing?")
	}
	rivalInfo.MainUser = true
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(rivalInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := s.syncRivalScoreLog(tx, rivalInfo); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
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

func (s *RivalInfoService) SyncRivalScoreLog(rivalInfo *entity.RivalInfo) error {
	if rivalInfo == nil {
		return fmt.Errorf("assert: SyncRivalScoreLog: rivalInfo == nil")
	}
	if rivalInfo.ID == 0 {
		return fmt.Errorf("assert: SyncRivalScoreLog: rivalInfo.ID corrupted")
	}
	log.Debug("[Service] calling RivalInfoService.SyncRivalScoreLog")
	if rivalInfo.ScoreLogPath == nil || *rivalInfo.ScoreLogPath == "" {
		return fmt.Errorf("cannot sync rival %s's score log: score log file path is empty", rivalInfo.Name)
	}
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := s.syncRivalScoreLog(tx, rivalInfo); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
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

func (s *RivalInfoService) FindRivalScoreLogByRivalId(ID uint) ([]entity.RivalScoreLog, error) {
	var logs []entity.RivalScoreLog
	if err := s.db.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *RivalInfoService) QueryUserInfoWithLevelLayeredDiffTableLampStatus(rivalID uint, headerID uint) (*dto.RivalInfoDto, error) {
	log.Debugf("[RivalInfoService] QueryUserInfoWithLevelLayeredDiffTableLampStatus: rivalID=%d, headerId=%d", rivalID, headerID)
	rivalInfo, err := s.FindRivalInfoByID(rivalID)
	if err != nil {
		return nil, err
	}
	header, err := s.diffTableService.QueryLevelLayeredDiffTableInfoById(headerID)
	if err != nil {
		return nil, err
	}
	logs, err := s.FindRivalScoreLogByRivalId(rivalID)
	if err != nil {
		return nil, err
	}

	sha256MaxLamp := make(map[string]int32)
	for _, log := range logs {
		if _, ok := sha256MaxLamp[log.Sha256]; !ok {
			sha256MaxLamp[log.Sha256] = log.Clear
		}
		sha256MaxLamp[log.Sha256] = max(sha256MaxLamp[log.Sha256], log.Clear)
	}

	ret := dto.NewRivalInfoDtoWithDiffTable(rivalInfo, header)
	for _, dataList := range ret.DiffTableHeader.LevelLayeredContents {
		for i, data := range dataList {
			if _, ok := sha256MaxLamp[data.Sha256]; ok {
				dataList[i].Lamp = int(sha256MaxLamp[data.Sha256])
			}
		}

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

// TODO: 移除掉该接口同步SongData数据的能力
func (s *RivalInfoService) syncRivalScoreLog(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
	// (1) load and sync scorelog.db
	if rivalInfo.ScoreLogPath == nil {
		return fmt.Errorf("assert: rival's scorelog path cannot be empty")
	}
	rawScoreLog, err := loadScoreLog(*rivalInfo.ScoreLogPath)
	if err != nil {
		return err
	}
	if err := syncScoreLog(tx, rawScoreLog, rivalInfo.ID); err != nil {
		return err
	}
	// (2) load and sync songdata.db
	rawSongData, err := loadSongData(*rivalInfo.SongDataPath)
	if err != nil {
		return err
	}
	if err := syncSongData(tx, rawSongData, rivalInfo.ID); err != nil {
		return err
	}
	// (3) generate rival tags
	if err := s.rivalTagService.syncRivalTagFromRawData(tx, rivalInfo.ID, rawScoreLog, rawSongData); err != nil {
		return err
	}
	// (4) update rival itself
	return tx.Model(rivalInfo).Updates(entity.RivalInfo{
		PlayCount: len(rawScoreLog),
	}).Error
}

func loadScoreLog(scoreLogPath string) ([]*entity.ScoreLog, error) {
	if scoreLogPath == "" {
		return nil, fmt.Errorf("assert: scorelog path cannot be empty")
	}
	log.Debugf("[RivalInfoService] Trying to read score log from %s", scoreLogPath)
	scoreLogDB, err := gorm.Open(sqlite.Open(scoreLogPath))
	if err != nil {
		return nil, err
	}
	scoreLogService := NewScoreLogService(scoreLogDB)
	rawScoreLog, n, err := scoreLogService.FindScoreLogList()
	if err != nil {
		return nil, err
	}
	log.Debugf("[RivalInfoService] Read %d logs from %s", n, scoreLogPath)
	return rawScoreLog, nil
}

func loadSongData(songDataPath string) ([]*entity.SongData, error) {
	if songDataPath == "" {
		return nil, fmt.Errorf("assert: songdata path cannot be empty")
	}
	songDataDB, err := gorm.Open(sqlite.Open(songDataPath))
	if err != nil {
		return nil, err
	}
	songDataService := NewSongDataService(songDataDB)
	rawSongData, n, err := songDataService.FindSongDataList()
	if err != nil {
		return nil, err
	}
	log.Infof("[RivalInfoService] Read %d song data from %s", n, songDataPath)
	return rawSongData, nil
}

func syncScoreLog(tx *gorm.DB, rawScoreLog []*entity.ScoreLog, rivalID uint) error {
	rivalScoreLog := make([]entity.RivalScoreLog, len(rawScoreLog))
	for i, rawLog := range rawScoreLog {
		rivalLog := entity.FromRawScoreLogToRivalScoreLog(rawLog)
		rivalLog.RivalId = rivalID
		rivalScoreLog[i] = rivalLog
	}
	if err := tx.Unscoped().Where("rival_id = ?", rivalID).Delete(&entity.RivalScoreLog{}).Error; err != nil {
		return err
	}

	if err := tx.CreateInBatches(&rivalScoreLog, 100).Error; err != nil {
		return err
	}
	return nil
}

func syncSongData(tx *gorm.DB, rawSongData []*entity.SongData, rivalID uint) error {
	rivalSongData := make([]entity.RivalSongData, len(rawSongData))
	for i, rawData := range rawSongData {
		rivalData := entity.FromRawSongDataToRivalSongData(rawData)
		rivalData.RivalId = rivalID
		rivalSongData[i] = rivalData
	}

	if err := tx.Unscoped().Where("rival_id = ?", rivalID).Delete(&entity.RivalSongData{}).Error; err != nil {
		return err
	}

	if err := tx.CreateInBatches(&rivalSongData, 100).Error; err != nil {
		return err
	}
	return nil
}

func queryMainUser(tx *gorm.DB) (*entity.RivalInfo, error) {
	var out entity.RivalInfo
	if err := tx.Where(&entity.RivalInfo{MainUser: true}).First(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
