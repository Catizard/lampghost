package service

import (
	"fmt"
	"os"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const READONLY_PARAMETER = "?open_mode=1"

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
	if rivalInfo.Name == "" {
		return fmt.Errorf("rival name cannot be empty")
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

// Wrapper method for SyncRivalData and IncrementalSyncRivalData, dispatching by config
func (s *RivalInfoService) SyncRivalDataByID(rivalID uint) error {
	if rivalInfo, err := s.FindRivalInfoByID(rivalID); err != nil {
		return err
	} else {
		conf, err := config.ReadConfig()
		if err != nil {
			return err
		}
		if conf.ForceFullyReload != 0 {
			log.Debug("[RivalInfoService] dispatched into fully reload")
			return s.SyncRivalData(rivalInfo)
		}
		log.Debug("[RivalInfoService] dispatched into incremental reload")
		return s.IncrementalSyncRivalData(rivalInfo)
	}
}

// Fully reload one rival's save files
func (s *RivalInfoService) SyncRivalData(rivalInfo *entity.RivalInfo) error {
	if rivalInfo == nil {
		return fmt.Errorf("assert: SyncRivalData: rivalInfo == nil")
	}
	if rivalInfo.ID == 0 {
		return fmt.Errorf("assert: SyncRivalData: rivalInfo.ID corrupted")
	}
	log.Debug("[Service] calling RivalInfoService.SyncRivalData")
	if rivalInfo.ScoreLogPath == nil || *rivalInfo.ScoreLogPath == "" {
		return fmt.Errorf("cannot sync rival %s's score log: score log file path is empty", rivalInfo.Name)
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := syncRivalData(tx, rivalInfo); err != nil {
			return err
		}
		if err := updateRivalPlayCount(tx, rivalInfo.ID); err != nil {
			return err
		}
		return nil
	})
}

// Extension to SyncRivalData, which only reloads part of the scorelog.db file
// More specifically, only reloads the log that is recorded after rival's last log
//
// Requirements:
//  1. rivalInfo's id > 0
//  2. rivalInfo's scorelog path must not be empty
//
// Special Cases:
// If no record belong to passed rival, fallback to fully reload
func (s *RivalInfoService) IncrementalSyncRivalData(rivalInfo *entity.RivalInfo) error {
	if rivalInfo == nil {
		return fmt.Errorf("incrementalSyncRivalData: rivalInfo cannot be nil")
	}
	if rivalInfo.ID == 0 {
		return fmt.Errorf("incrementalSyncRivalData: rivalInfo.ID should > 0")
	}
	if rivalInfo.ScoreLogPath == nil || *rivalInfo.ScoreLogPath == "" {
		return fmt.Errorf("incrementalSyncRivalData: rivalInfo.ScoreLogPath cannot be empty")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := incrementalSyncRivalData(tx, rivalInfo); err != nil {
			return err
		}
		if err := updateRivalPlayCount(tx, rivalInfo.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *RivalInfoService) DelRivalInfo(ID uint) error {
	mainUser, err := queryMainUser(s.db)
	if err != nil {
		return err
	}
	if ID == mainUser.ID {
		return fmt.Errorf("DelRivalInfo: cannot delete main user")
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		var candidate entity.RivalInfo
		if err := tx.First(&candidate, ID).Error; err != nil {
			return err
		}
		// RivalInfo
		if err := tx.Delete(&entity.RivalInfo{}, candidate.ID).Error; err != nil {
			return err
		}
		// RivalScoreLog
		if err := tx.Unscoped().Where("rival_id = ?", candidate.ID).Delete(&entity.RivalScoreLog{}).Error; err != nil {
			return err
		}
		// RivalSongData
		if err := tx.Unscoped().Where("rival_id = ?", candidate.ID).Delete(&entity.RivalSongData{}).Error; err != nil {
			return err
		}
		// RivalTag
		if err := tx.Unscoped().Where("rival_id = ?", candidate.ID).Delete(entity.RivalTag{}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *RivalInfoService) QueryUserPlayCountInYear(ID uint, yearNum string) ([]int, error) {
	logs, _, err := findRivalScoreLogList(s.db, &vo.RivalScoreLogVo{
		RivalId:     ID,
		SpecifyYear: &yearNum,
	})
	if err != nil {
		return nil, err
	}
	ret := make([]int, 12)
	for i := range ret {
		ret[i] = 0
	}
	for _, playLog := range logs {
		ret[playLog.RecordTime.Month()-1]++
	}
	return ret, nil
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
	sha256MaxLamp, err := findRivalMaximumClearScoreLogSha256Map(s.db, &vo.RivalScoreLogVo{
		RivalId: rivalID,
	})
	if err != nil {
		return nil, err
	}

	ret := dto.NewRivalInfoDtoWithDiffTable(rivalInfo, header)
	for _, dataList := range ret.DiffTableHeader.LevelLayeredContents {
		for i, data := range dataList {
			ret.DiffTableHeader.SongCount++
			if _, ok := sha256MaxLamp[data.Sha256]; ok {
				dataList[i].Lamp = int(sha256MaxLamp[data.Sha256][0].Clear)
				ret.DiffTableHeader.LampCount[dataList[i].Lamp]++
			}
		}
	}
	return ret, nil
}

func (s *RivalInfoService) QueryRivalPlayedYears(rivalID uint) ([]int, int, error) {
	var years []int
	if err := s.db.Model(&entity.RivalScoreLog{}).Where("rival_id = ?", rivalID).Select("distinct STRFTIME('%Y', record_time)").Find(&years).Error; err != nil {
		return nil, 0, err
	}
	return years, len(years), nil
}

func (s *RivalInfoService) QueryMainUser() (*entity.RivalInfo, error) {
	return queryMainUser(s.db)
}

func (s *RivalInfoService) UpdateRivalInfo(rivalInfo *vo.RivalInfoVo) error {
	if rivalInfo == nil {
		return fmt.Errorf("UpdateRivalInfo: rivalInfo cannot be nil")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		return updateRivalInfo(tx, rivalInfo.Entity())
	})
}

// Add one rival
//
// Don't call incremental sync in this method
func addRivalInfo(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
	if err := tx.Create(rivalInfo).Error; err != nil {
		return err
	}
	if err := syncRivalData(tx, rivalInfo); err != nil {
		return err
	}
	if err := updateRivalPlayCount(tx, rivalInfo.ID); err != nil {
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
// WARN: This function wouldn't update rival's play count, it's not this function's purpose
//
// For now, theses files would be reloaded
//  1. songdata.db
//  2. scorelog.db
//
// And these tables' data would be re-generated
//  1. rival_score_log
//  2. rival_song_data
//  3. rival_tag
//
// Specially, songdata.db file path is not provide, this method wouldn't return any error but skip.
// This hack is because we haven't implmeneted the correct seperate songdata.db for different rivals
// TODO: This hack also breaks some related data build
func syncRivalData(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
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
	// (2) load and sync songdata.db if if is provided
	if rivalInfo.SongDataPath != nil && *rivalInfo.SongDataPath != "" {
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
		// invalidate default song cache since we have rebuilt the `rival_song_data` table
		expireDefaultCache()
	} else {
		// (3) generate rival tags
		// Since we don't have rawSongData here, we need to build tags by using the default song data cache
		if err := syncRivalTag(tx, rivalInfo.ID); err != nil {
			return err
		}
	}
	return nil
}

// Incrementally reload one rival's save files
//
// WARN: This function wouldn't update rival's play count, it's not this function's purpose
//
// For now, theses files would be reloaded
//  1. scorelog.db
//
// And these tables' data would be re-generated
//  1. rival_score_log
//  2. rival_tag
//
// NOTE: This function wouldn't read `songdata.db` file, therefore there is no need to invalidate the
// default song cache
func incrementalSyncRivalData(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
	lastRivalScoreLog, err := findLastRivalScoreLogList(tx, &vo.RivalScoreLogVo{RivalId: rivalInfo.ID})
	if err != nil {
		return err
	}
	// Fallback to fully reload
	if lastRivalScoreLog.ID == 0 {
		return syncRivalData(tx, rivalInfo)
	}
	maximumRecordTimestamp := lastRivalScoreLog.RecordTime.Unix()
	rawScoreLog, err := loadScoreLog(*rivalInfo.ScoreLogPath, &maximumRecordTimestamp)
	if err != nil {
		return err
	}
	if err := appendScoreLog(tx, rawScoreLog, rivalInfo.ID); err != nil {
		return err
	}
	return syncRivalTag(tx, rivalInfo.ID)
}

// Read one `scorelog.db` file into memory
func loadScoreLog(scoreLogPath string, maximumTimestamp *int64) ([]*entity.ScoreLog, error) {
	if scoreLogPath == "" {
		return nil, fmt.Errorf("assert: scorelog path cannot be empty")
	}
	log.Debugf("[RivalInfoService] Trying to read score log from %s", scoreLogPath)
	if err := verifyLocalDatabaseFilePath(scoreLogPath); err != nil {
		return nil, err
	}
	dsn := scoreLogPath + READONLY_PARAMETER
	log.Debugf("[RivalInfoService] Full scorelog.db dsn: %s", dsn)
	scoreLogDB, err := gorm.Open(sqlite.Open(dsn))
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

// Read one `songdata.db` file into memory
func loadSongData(songDataPath string) ([]*entity.SongData, error) {
	if songDataPath == "" {
		return nil, fmt.Errorf("assert: songdata path cannot be empty")
	}
	log.Debugf("[RivalInfoService] Trying to read song data from %s", songDataPath)
	if err := verifyLocalDatabaseFilePath(songDataPath); err != nil {
		return nil, err
	}
	dsn := songDataPath + READONLY_PARAMETER
	songDataDB, err := gorm.Open(sqlite.Open(dsn))
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

// Helper function for validating local databasement file path
func verifyLocalDatabaseFilePath(filePath string) error {
	if stat, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("assert: no file exists at %s", filePath)
		}
		return fmt.Errorf("assert: cannot stat file at %s", filePath)
	} else if stat.IsDir() {
		return fmt.Errorf("assert: file path %s is a directory, not an valid database file", filePath)
	}
	return nil
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

// Update one rival info
//
// Fully reload save files if one of theses file paths not nil and has been changed:
//
//	1.ScoreLogPath
//	2.SongDataPath
//
// Special Requirements:
//  1. when updating main user, songdata.db file path cannot be empty
//
// Warning: You should never be able to edit `MainUser` by using update interface
func updateRivalInfo(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
	if rivalInfo.ID == 0 {
		return fmt.Errorf("updateRivalInfo: rivalInfo.ID cannot be 0")
	}

	var prev entity.RivalInfo
	if err := tx.First(&prev, rivalInfo.ID).Error; err != nil {
		return err
	}

	if prev.MainUser && (rivalInfo.SongDataPath == nil || *rivalInfo.SongDataPath == "") {
		return fmt.Errorf("updateRivalInfo: SongDataPath cannot be empty when updateing main user")
	}

	// TODO: we haven't implemented the correct seperate songdata file feature
	if !prev.MainUser && (rivalInfo.SongDataPath != nil && *rivalInfo.SongDataPath != "") {
		return fmt.Errorf("updateRivalInfo: cannot provide songdata path for a non main user")
	}

	shouldFullyReload := false
	if rivalInfo.ScoreLogPath != nil && *rivalInfo.ScoreLogPath != *prev.ScoreLogPath {
		shouldFullyReload = true
	}
	if rivalInfo.SongDataPath != nil && *rivalInfo.SongDataPath != *prev.SongDataPath {
		shouldFullyReload = true
	}

	if shouldFullyReload {
		if err := syncRivalData(tx, rivalInfo); err != nil {
			return err
		}

		// Since we have fully reloaded save files, we need also update play count here
		pc, err := selectRivalScoreLogCount(tx, &vo.RivalScoreLogVo{RivalId: rivalInfo.ID})
		if err != nil {
			return err
		}
		rivalInfo.PlayCount = int(pc)
	}

	return tx.Updates(rivalInfo).Error
}

// Simple helper function for updating one rival's play count field
func updateRivalPlayCount(tx *gorm.DB, rivalID uint) error {
	if rivalID == 0 {
		return fmt.Errorf("updateRivalPlayCount: rivalID cannot be 0")
	}

	pc, err := selectRivalScoreLogCount(tx, &vo.RivalScoreLogVo{RivalId: rivalID})
	if err != nil {
		return err
	}
	return tx.Model(&entity.RivalInfo{Model: gorm.Model{ID: rivalID}}).UpdateColumn("play_count", pc).Error
}
