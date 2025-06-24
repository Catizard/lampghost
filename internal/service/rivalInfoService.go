package service

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"github.com/mitchellh/mapstructure"
	"github.com/rotisserie/eris"
	. "github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

const READONLY_PARAMETER = "?open_mode=1"

type RivalInfoService struct {
	db  *gorm.DB
	ctx context.Context
}

func NewRivalInfoService(db *gorm.DB) *RivalInfoService {
	return &RivalInfoService{
		db: db,
	}
}

func (s *RivalInfoService) InjectContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *RivalInfoService) InitializeMainUser(rivalInfo *vo.InitializeRivalInfoVo) error {
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
	// TODO: Wtf is this, does it do anything?
	config, err := config.ReadConfig()
	if err != nil {
		return err
	}
	config.WriteConfig()
	insertRivalInfo := rivalInfo.Into()
	insertRivalInfo.MainUser = true
	// Prechecks
	if insertRivalInfo.SongDataPath == nil || *insertRivalInfo.SongDataPath == "" {
		return eris.New("songdata.db path cannot be empty")
	}
	if insertRivalInfo.ScoreLogPath == nil || *insertRivalInfo.ScoreLogPath == "" {
		return eris.New("scorelog.db path cannot be empty")
	}
	if insertRivalInfo.ScoreDataLogPath == nil || *insertRivalInfo.ScoreDataLogPath == "" {
		return eris.New("scoredatalog.db path cannot be empty")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		return addRivalInfo(tx, insertRivalInfo)
	})
}

// Open a choose directory dialog and verify whether it's a valid
// beatoraja directory or not. Returning error if it's not.
//
// A valid beatoraja directory is defined as follow:
//  1. contains a 'songdata.db'
//  2. contains a directory called 'player'
//  3. 'playder' directory is not empty
//
// Returns directory names which is located under 'player' and
// beatoraja directory path chosen by user
func (s *RivalInfoService) ChooseBeatorajaDirectory() (*dto.BeatorajaDirectoryMeta, error) {
	dir, err := runtime.OpenDirectoryDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "Choose Beatoraja Directory",
	})
	if err != nil {
		return nil, err
	}
	if dir == "" {
		return nil, eris.New("choose directory: directory cannot be empty")
	}

	// 1. contains a 'songdata.db'
	songdataPath := path.Join(dir, "songdata.db")
	if err := database.VerifyLocalDatabaseFilePath(songdataPath); err != nil {
		return nil, eris.Wrap(err, "not a valid songdata.db file path")
	}
	// 2. contains a directory called 'player'
	playerDirectoryPath := path.Join(dir, "player")
	if stat, err := os.Stat(playerDirectoryPath); err != nil {
		return nil, eris.Wrap(err, "cannot stat player directory")
	} else if !stat.Mode().IsDir() {
		return nil, eris.Wrap(err, "player is not a directory")
	} else {
		// Okay, it's a valid directory
		entries, err := os.ReadDir(playerDirectoryPath)
		if err != nil {
			return nil, eris.Wrap(err, "cannot read player directory")
		}
		// 3. 'player' directory is not empty
		if len(entries) == 0 {
			return nil, eris.New("player directory is empty")
		}
		possibleSaveDirectories := make([]string, 0)
		for _, entry := range entries {
			if entry.IsDir() {
				possibleSaveDirectories = append(possibleSaveDirectories, entry.Name())
			}
		}
		return &dto.BeatorajaDirectoryMeta{
			BeatorajaDirectoryPath: dir,
			PlayerDirectories:      possibleSaveDirectories,
		}, nil
	}
}

// Add one non-main user
//
// NOTE: Currently, we don't support seperate songdata.db feature, therefore,
// this function would return error if user trying to add a non-main user with
// a seperate songdata.db file path
func (s *RivalInfoService) AddRivalInfo(rivalInfo *vo.RivalInfoVo) error {
	if rivalInfo == nil {
		return eris.Errorf("add rival: rivalInfo cannot be nil")
	}
	if rivalInfo.ScoreLogPath == nil || *rivalInfo.ScoreLogPath == "" {
		return eris.Errorf("add rival: scorelog.db path cannot be empty")
	}
	if rivalInfo.SongDataPath != nil && *rivalInfo.SongDataPath != "" {
		return eris.Errorf("add rival: seperate songdata.db is not supported currently")
	}
	if rivalInfo.Name == "" {
		return fmt.Errorf("rival name cannot be empty")
	}
	// No, you can never add a main user by using this inteface
	rivalInfo.MainUser = false
	return s.db.Transaction(func(tx *gorm.DB) error {
		return addRivalInfo(tx, rivalInfo.Entity())
	})
}

func (s *RivalInfoService) FindRivalInfoList(filter *vo.RivalInfoVo) ([]*dto.RivalInfoDto, int, error) {
	return findRivalInfoList(s.db, filter)
}

func (s *RivalInfoService) FindRivalInfoByID(rivalID uint) (*entity.RivalInfo, error) {
	out := entity.RivalInfo{}
	if err := s.db.First(&out, rivalID).Error; err != nil {
		log.Debugf("[RivalInfoService] FindRivalInfoByID with ID=%d failed: %v\n", rivalID, err)
		return nil, err
	}
	return &out, nil
}

// Reload one rival's data, using different strategy based on fullyReload's value
//
//  1. fullyReload == 0: Incrementallay update `scorelog.db` and `scoredatalog.db` file, nothing to do with `songdata.db`
//  2. fullyReload == 1: Fully reload every files, for now it's `scorelog.db`, `songdata.db` and `scoredatalog.db`
func (s *RivalInfoService) ReloadRivalData(rivalID uint, fullyReload bool) error {
	if rivalInfo, err := s.FindRivalInfoByID(rivalID); err != nil {
		return eris.Wrapf(err, "cannot find user")
	} else {
		if fullyReload {
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
	header, err := queryLevelLayeredDiffTableInfoById(s.db, headerID)
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

// Warning: You should never be able to edit `MainUser` by using update interface
func (s *RivalInfoService) UpdateRivalInfo(rivalInfo *vo.RivalInfoVo) error {
	if rivalInfo == nil {
		return eris.Errorf("rivalInfo cannot be nil")
	}
	if rivalInfo.ID == 0 {
		return eris.Errorf("ID cannot be 0")
	}
	// Unset crucial fields
	rivalInfo.PlayCount = 0
	rivalInfo.MainUser = false
	log.Debug("before opening transaction")
	return s.db.Transaction(func(tx *gorm.DB) error {
		return updateRivalInfo(tx, rivalInfo.Entity())
	})
}

// Special variant of 'UpdateRivalInfo' function, for updating 'Reverse Import' feature fields
//
// Frontend shows the fundemental fields (save files location, name...) and reverse import fields
// (lock tag id, reverse import flag) in different components with different forms.
// And it's painful for updating plain value fields when using gorm :(
// Therefore, spliting the update functions seems feasible way currently
func (s *RivalInfoService) UpdateRivalReverseImportInfo(rivalInfo *vo.RivalInfoVo) error {
	if rivalInfo == nil {
		return eris.Errorf("UpdateRivalReverseImportInfo: rivalInfo cannot be nil")
	}
	if rivalInfo.ID == 0 {
		return eris.Errorf("UpdateRivalReverseImportInfo: ID cannot be 0")
	}
	// Only preserve necessary fields
	updateParam := vo.RivalInfoVo{
		Model: gorm.Model{
			ID: rivalInfo.ID,
		},
		ReverseImport: rivalInfo.ReverseImport,
		LockTagID:     rivalInfo.LockTagID,
	}
	if err := updateRivalInfoPlainFields(s.db, *updateParam.Entity()); err != nil {
		return eris.Wrap(err, "update rival info plain fields")
	}
	return nil
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

func findRivalInfoList(tx *gorm.DB, filter *vo.RivalInfoVo) ([]*dto.RivalInfoDto, int, error) {
	var out []*dto.RivalInfoDto
	fields := `
		rival_info.*,
		rival_tag.tag_name
	`
	moved := tx.Select(fields).Model(&entity.RivalInfo{}).Scopes(scopeRivalInfoFilter(filter))
	moved = moved.Joins(`left join rival_tag on rival_info.lock_tag_id = rival_tag.id`)

	if err := moved.Find(&out).Error; err != nil {
		return nil, 0, err
	}

	return out, len(out), nil
}

func selectRivalInfoCount(tx *gorm.DB, filter *vo.RivalInfoVo) (int64, error) {
	var count int64
	if err := tx.Model(&entity.RivalInfo{}).Scopes(scopeRivalInfoFilter(filter)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Cmmon query scope for vo.RivalInfoVo
func scopeRivalInfoFilter(filter *vo.RivalInfoVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Add extra filter here
		if filter.IgnoreMainUser {
			moved.Where("ID != 1")
		}
		return moved
	}
}

// Fully reload one rival's saves files
//
// WARN: This function wouldn't update rival's play count, it's not this
// function's purpose
//
// For now, theses files would be reloaded
//  1. scorelog.db
//  2. songdata.db (if provides)
//  3. scoredatalog.db (if provides)
//
// And these tables' data would be re-generated
//  1. rival_score_log
//  2. rival_song_data (if songdata.db provides)
//  3. rival_score_data_log (if scoredatalog.db provides)
//  4. rival_tag
//
// NOTE: The way we re-generate rival's tags is based on whether songdata.db
// file exists or not, see below implementation for details.
//
// Specially, songdata.db or scoredatalog.db file path is not provide, this
// method wouldn't return any error but skip.
// This hack is because we haven't implmeneted the correct seperate
// songdata.db for different rivals and scoredatalog.db is optional.
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
	// (2) load and sync songdata.db if it's provided
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
	// (4) load and sync scoredatalog.db if it's provided
	if rivalInfo.ScoreDataLogPath != nil && *rivalInfo.ScoreDataLogPath != "" {
		rawScoreDataLog, err := loadScoreDataLog(*rivalInfo.ScoreDataLogPath, nil)
		if err != nil {
			return err
		}
		if err := syncScoreDataLog(tx, rawScoreDataLog, rivalInfo.ID); err != nil {
			return err
		}
	}

	return nil
}

// Incrementally reload one rival's save files
//
// WARN: This function wouldn't update rival's play count, it's not
// this function's purpose
//
// For now, theses files would be reloaded
//  1. scorelog.db
//  2. scoredatalog.db (if provides)
//
// And these tables' data would be updated
//  1. rival_score_log (incrementally added)
//  2. rival_score_data_log (incrementally added)
//  3. rival_tag (keep to old data as much as possible)
//
// NOTE: This function wouldn't read `songdata.db` file, therefore there
// is no need to invalidate the default song cache
//
// NOTE: This function would degenerate to fully reload if:
//  1. rival_score_log is empty
//  2. rival_score_data_log is empty (This wouldn't be a problem since songdata.db could only be imported by main user currently)
func incrementalSyncRivalData(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
	lastRivalScoreLog, err := findLastRivalScoreLog(tx, &vo.RivalScoreLogVo{RivalId: rivalInfo.ID})
	if err != nil {
		return err
	}
	lastRivalScoreDataLog, err := findLastRivalScoreDataLog(tx, &vo.RivalScoreDataLogVo{RivalId: rivalInfo.ID})
	if err != nil {
		return err
	}
	// Fallback to fully reload
	if lastRivalScoreLog.ID == 0 || lastRivalScoreDataLog.ID == 0 {
		return syncRivalData(tx, rivalInfo)
	}
	scoreLogMaximumRecordTimestamp := lastRivalScoreLog.RecordTime.Unix()
	rawScoreLog, err := loadScoreLog(*rivalInfo.ScoreLogPath, &scoreLogMaximumRecordTimestamp)
	if err != nil {
		return err
	}
	if err := appendScoreLog(tx, rawScoreLog, rivalInfo.ID); err != nil {
		return err
	}
	scoreDataLogMaximumRecordTimestamp := lastRivalScoreDataLog.RecordTime.Unix()
	rawScoreDataLog, err := loadScoreDataLog(*rivalInfo.ScoreDataLogPath, &scoreDataLogMaximumRecordTimestamp)
	if err != nil {
		return err
	}
	if err := appendScoreDataLog(tx, rawScoreDataLog, rivalInfo.ID); err != nil {
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
	if err := database.VerifyLocalDatabaseFilePath(scoreLogPath); err != nil {
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

// Read all contents from `songdata.db` file into memory
func loadSongData(songDataPath string) ([]*entity.SongData, error) {
	if songDataPath == "" {
		return nil, fmt.Errorf("assert: songdata path cannot be empty")
	}
	log.Debugf("[RivalInfoService] Trying to read song data from %s", songDataPath)
	if err := database.VerifyLocalDatabaseFilePath(songDataPath); err != nil {
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

// Read all contents from `scoredatalog.db` file into memory
func loadScoreDataLog(scoreDataLogPath string, maximumTimestamp *int64) ([]*entity.ScoreDataLog, error) {
	if scoreDataLogPath == "" {
		return nil, eris.New("load: scoredatalog.db file path cannot be empty")
	}
	log.Debugf("[RivalInfoService] Trying to read log from %s", scoreDataLogPath)
	if err := database.VerifyLocalDatabaseFilePath(scoreDataLogPath); err != nil {
		return nil, err
	}
	dsn := scoreDataLogPath + READONLY_PARAMETER
	scoreDataLogDB, err := gorm.Open(sqlite.Open(dsn))
	if err != nil {
		return nil, eris.Wrap(err, "failed to open scoredatalog.db")
	}
	scoreDataLogService := NewScoreDataLogService(scoreDataLogDB)
	rawScoreDataLog, n, err := scoreDataLogService.FindScoreDataLogList(maximumTimestamp)
	if err != nil {
		return nil, eris.Wrap(err, "load: query from scoredatalog.db failed")
	}
	log.Debugf("[RivalInfoService] Read %d logs from %s", n, scoreDataLogPath)
	return rawScoreDataLog, nil
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

// Similar to syncScoreDataLog but not delete any old content, only append new logs
func appendScoreDataLog(tx *gorm.DB, rawScoreDataLog []*entity.ScoreDataLog, rivalID uint) error {
	newRivalScoreDataLogs := Map(rawScoreDataLog, func(rawLog *entity.ScoreDataLog, _ int) *entity.RivalScoreDataLog {
		ret := entity.FromRawScoreDataLogToRivalScoreDataLog(rawLog)
		ret.RivalId = rivalID
		return &ret
	})
	return tx.Model(&entity.RivalScoreDataLog{}).CreateInBatches(newRivalScoreDataLogs, DEFAULT_BATCH_SIZE).Error
}

// Fully delete all content from rival_score_log and rebuild them by rawScoreLog
func syncSongData(tx *gorm.DB, rawSongData []*entity.SongData, rivalID uint) error {
	rivalSongData := Map(rawSongData, func(rawData *entity.SongData, _ int) *entity.RivalSongData {
		ret := entity.FromRawSongDataToRivalSongData(rawData)
		ret.RivalId = rivalID
		return &ret
	})

	if err := tx.Unscoped().Where("rival_id = ?", rivalID).Delete(&entity.RivalSongData{}).Error; err != nil {
		return err
	}

	return tx.CreateInBatches(&rivalSongData, DEFAULT_BATCH_SIZE).Error
}

// Fully delete all content from rival_score_data_log and rebuild them by rawScoreDataLog
func syncScoreDataLog(tx *gorm.DB, rawScoreDataLog []*entity.ScoreDataLog, rivalID uint) error {
	rivalScoreDataLog := Map(rawScoreDataLog, func(rawLog *entity.ScoreDataLog, _ int) *entity.RivalScoreDataLog {
		ret := entity.FromRawScoreDataLogToRivalScoreDataLog(rawLog)
		ret.RivalId = rivalID
		return &ret
	})

	if err := tx.Unscoped().Where("rival_id = ?", rivalID).Delete(&entity.RivalScoreDataLog{}).Error; err != nil {
		return err
	}

	return tx.CreateInBatches(&rivalScoreDataLog, DEFAULT_BATCH_SIZE).Error
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
//  1. ScoreLogPath
//  2. SongDataPath
//  3. ScoreDataLogPath
//
// Special Requirements:
//  1. when updating main user, songdata.db file path cannot be empty
func updateRivalInfo(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
	var prev entity.RivalInfo
	if err := tx.Debug().First(&prev, rivalInfo.ID).Error; err != nil {
		return err
	}

	if prev.MainUser && (rivalInfo.SongDataPath == nil || *rivalInfo.SongDataPath == "") {
		return eris.Errorf("updateRivalInfo: SongDataPath cannot be empty when updateing main user")
	}

	// TODO: we haven't implemented the correct seperate songdata file feature
	if !prev.MainUser && (rivalInfo.SongDataPath != nil && *rivalInfo.SongDataPath != "") {
		return eris.Errorf("updateRivalInfo: cannot provide songdata path for a non main user")
	}

	shouldFullyReload := false
	if rivalInfo.ScoreLogPath != nil && *rivalInfo.ScoreLogPath != *prev.ScoreLogPath {
		shouldFullyReload = true
	}
	if rivalInfo.SongDataPath != nil && *rivalInfo.SongDataPath != *prev.SongDataPath {
		shouldFullyReload = true
	}
	if rivalInfo.ScoreDataLogPath != nil && (prev.ScoreDataLogPath == nil || *prev.ScoreDataLogPath != *rivalInfo.ScoreDataLogPath) {
		shouldFullyReload = true
	}

	if shouldFullyReload {
		log.Debugf("[RivalInfoService] trying to fully update rival info")
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

// Update one rival's plain fields, no special treatment is being done
// The caller side must know what they are doing
func updateRivalInfoPlainFields(tx *gorm.DB, rivalInfo entity.RivalInfo) error {
	updates := new(map[string]any)
	if err := mapstructure.Decode(rivalInfo, updates); err != nil {
		return eris.Wrap(err, "cannot map rival info into map")
	}
	return eris.Wrap(tx.Model(&rivalInfo).Updates(updates).Error, "cannot update rival_info")
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
	return tx.Model(&entity.RivalInfo{Model: gorm.Model{ID: rivalID}}).Update("play_count", pc).Error
}
