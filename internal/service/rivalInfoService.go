package service

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const DEFAULT_BATCH_SIZE = 100

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

func (s *RivalInfoService) InitializeMainUser(rivalInfo *vo.RivalInfoVo) error {
	if rivalInfo.SongDataPath == nil || *rivalInfo.SongDataPath == "" {
		return fmt.Errorf("songdata.db path cannot be empty")
	}
	if rivalInfo.ScoreLogPath == nil || *rivalInfo.ScoreLogPath == "" {
		return fmt.Errorf("scorelog.db path cannot be empty")
	}
	if rivalInfo.Locale != nil && *rivalInfo.Locale != "" {
		conf, err := config.ReadConfig()
		if err != nil {
			return err
		}
		conf.Locale = *rivalInfo.Locale
		if err := conf.WriteConfig(); err != nil {
			return err
		}
	}
	mainUserCount, err := selectRivalInfoCount(s.db, &vo.RivalInfoVo{MainUser: true})
	if err != nil {
		return err
	}
	if mainUserCount > 0 {
		return fmt.Errorf("cannot have two main user, what are you doing?")
	}
	// Initialize the config
	config, err := config.ReadConfig()
	if err != nil {
		return err
	}
	config.UserName = rivalInfo.Name
	config.SongdataFilePath = *rivalInfo.SongDataPath
	config.ScorelogFilePath = *rivalInfo.ScoreLogPath
	config.WriteConfig()
	rivalInfo.MainUser = true
	return s.db.Transaction(func(tx *gorm.DB) error {
		return addRivalInfo(tx, rivalInfo.Entity())
	})
}

func (s *RivalInfoService) AddRivalInfo(rivalInfo *entity.RivalInfo) error {
	if rivalInfo == nil {
		return fmt.Errorf("AddRivalInfo: rivalInfo cannot be nil")
	}
	if rivalInfo.ScoreLogPath == nil || *rivalInfo.ScoreLogPath == "" {
		return fmt.Errorf("scorelog.db path cannot be empty")
	}
	// No, you can never add a main user by using this inteface
	rivalInfo.MainUser = false
	return s.db.Transaction(func(tx *gorm.DB) error {
		return addRivalInfo(tx, rivalInfo)
	})
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

// Fully reload one rival's save files
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
	return s.db.Transaction(func(tx *gorm.DB) error {
		return syncRivalScoreLog(tx, rivalInfo)
	})
}

// Extension to syncRivalScoreLog, which only reloads part of the scorelog.db file
// More specifically, only reloads the log that is recorded after rival's last log
//
// Requirements:
// 1) rivalInfo's id > 0
// 2) rivalInfo's scorelog path must not be empty
//
// Special Cases:
// If no record belong to passed rival, fallback to fully reload
func (s *RivalInfoService) IncrementalSyncRivalScoreLog(rivalInfo *entity.RivalInfo) error {
	if rivalInfo == nil {
		return fmt.Errorf("incrementalSyncRivalScoreLog: rivalInfo cannot be nil")
	}
	if rivalInfo.ID == 0 {
		return fmt.Errorf("incrementalSyncRivalScoreLog: rivalInfo.ID should > 0")
	}
	if rivalInfo.ScoreLogPath == nil || *rivalInfo.ScoreLogPath == "" {
		return fmt.Errorf("incrementalSyncRivalScoreLog: rivalInfo.ScoreLogPath cannot be empty")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		return incrementalSyncRivalScoreLog(tx, rivalInfo)
	})
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
		ret[playLog.RecordTime.Month()-1]++
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
		conf, err := config.ReadConfig()
		if err != nil {
			return err
		}
		if conf.ForceFullyReload != 0 {
			log.Debug("[RivalInfoService] dispatched into fully reload")
			return s.SyncRivalScoreLog(rivalInfo)
		}
		log.Debug("[RivalInfoService] dispatched into incremental reload")
		return s.IncrementalSyncRivalScoreLog(rivalInfo)
	}
}

func (s *RivalInfoService) QueryMainUser() (*entity.RivalInfo, error) {
	var out entity.RivalInfo
	if err := s.db.Where(&entity.RivalInfo{MainUser: true}).First(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func addRivalInfo(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
	if err := tx.Create(rivalInfo).Error; err != nil {
		return err
	}
	if err := syncRivalScoreLog(tx, rivalInfo); err != nil {
		return err
	}
	return nil
}

func selectRivalInfoCount(tx *gorm.DB, filter *vo.RivalInfoVo) (int64, error) {
	querying := tx.Model(&entity.RivalInfo{})
	if filter != nil {
		querying = querying.Where(filter.Entity())
	}
	var count int64
	if err := querying.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Fully reload one rival's saves files
//
// Warning: this method actually reload not only scorelog.db but all save files
// may need a rename in the future
func syncRivalScoreLog(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
	// (1) load and sync scorelog.db
	if rivalInfo.ScoreLogPath == nil {
		return fmt.Errorf("assert: rival's scorelog path cannot be empty")
	}
	rawScoreLog, err := loadScoreLog(*rivalInfo.ScoreLogPath, nil)
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
	if err := syncRivalTagFromRawData(tx, rivalInfo.ID, rawScoreLog, rawSongData); err != nil {
		return err
	}
	// (4) update rival itself
	return tx.Model(rivalInfo).Updates(entity.RivalInfo{
		PlayCount: len(rawScoreLog),
	}).Error
}

func incrementalSyncRivalScoreLog(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
	lastRivalScoreLog, err := findLastRivalScoreLogList(tx, &vo.RivalScoreLogVo{RivalId: rivalInfo.ID})
	if err != nil {
		return err
	}
	// Fallback to fully reload
	if lastRivalScoreLog.ID == 0 {
		return syncRivalScoreLog(tx, rivalInfo)
	}
	maximumRecordTimestamp := lastRivalScoreLog.RecordTime.Unix()
	rawScoreLog, err := loadScoreLog(*rivalInfo.ScoreLogPath, &maximumRecordTimestamp)
	if err != nil {
		return err
	}
	if err := appendScoreLog(tx, rawScoreLog, rivalInfo.ID); err != nil {
		return err
	}
	syncRivalTag(tx, rivalInfo.ID)
	// TODO: update rival's playcount here
	return nil
}

func loadScoreLog(scoreLogPath string, maximumTimestamp *int64) ([]*entity.ScoreLog, error) {
	if scoreLogPath == "" {
		return nil, fmt.Errorf("assert: scorelog path cannot be empty")
	}
	log.Debugf("[RivalInfoService] Trying to read score log from %s", scoreLogPath)
	scoreLogDB, err := gorm.Open(sqlite.Open(scoreLogPath))
	if err != nil {
		return nil, err
	}
	scoreLogService := NewScoreLogService(scoreLogDB)
	rawScoreLog, n, err := scoreLogService.FindScoreLogList(maximumTimestamp)
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

// Fully delete all content from rival_score_log and rebuild them by rawScoreLog
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

	if err := tx.CreateInBatches(&rivalScoreLog, DEFAULT_BATCH_SIZE).Error; err != nil {
		return err
	}
	return nil
}

// Similar to syncScoreLog but not delete any old content, only append new logs
func appendScoreLog(tx *gorm.DB, rawScoreLog []*entity.ScoreLog, rivalID uint) error {
	newRivalScorelogs := make([]*entity.RivalScoreLog, len(rawScoreLog))
	for i, rawLog := range rawScoreLog {
		rivalLog := entity.FromRawScoreLogToRivalScoreLog(rawLog)
		rivalLog.RivalId = rivalID
		newRivalScorelogs[i] = &rivalLog
	}
	return tx.Model(&entity.RivalScoreLog{}).CreateInBatches(newRivalScorelogs, DEFAULT_BATCH_SIZE).Error
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

	return tx.CreateInBatches(&rivalSongData, DEFAULT_BATCH_SIZE).Error
}

func queryMainUser(tx *gorm.DB) (*entity.RivalInfo, error) {
	var out entity.RivalInfo
	if err := tx.Where(&entity.RivalInfo{MainUser: true}).First(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
