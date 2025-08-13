package service

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"github.com/rotisserie/eris"
	. "github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

const READONLY_PARAMETER = "?open_mode=1"

type RivalInfoService struct {
	db             *gorm.DB
	monitorService *MonitorService
	ctx            context.Context
	syncChan       <-chan any
}

func NewRivalInfoService(db *gorm.DB, monitorService *MonitorService, syncChan <-chan any) *RivalInfoService {
	ret := &RivalInfoService{
		db:             db,
		monitorService: monitorService,
		syncChan:       syncChan,
	}
	go ret.listen()
	return ret
}

func (s *RivalInfoService) listen() {
	for {
		<-s.syncChan
		runtime.EventsEmit(s.ctx, "global:notify", dto.NotificationDto{
			Type:    "info",
			Content: "File change detected, trying to auto-reload save files",
		})
		// TODO: Magical main user id=1
		if err := s.ReloadRivalData(1, false); err != nil {
			log.Errorf("failed to auto-reload: %s", err)
			runtime.EventsEmit(s.ctx, "global:notify", dto.NotificationDto{
				Type:    "error",
				Content: fmt.Sprintf("Failed to auto-reload: %s", err),
			})
		} else {
			runtime.EventsEmit(s.ctx, "global:refresh")
			// HACK: Make the message as visible as possible
			go func() {
				time.Sleep(1 * time.Second)
				runtime.EventsEmit(s.ctx, "global:notify", dto.NotificationDto{
					Type:    "success",
					Content: "Successfully auto-reload save files",
				})
			}()
		}
	}
}

func (s *RivalInfoService) InjectContext(ctx context.Context) {
	s.ctx = ctx
}

// TODO: I think it's better to change the signature to InitializeMainUser(conf, rivalInfo)
func (s *RivalInfoService) InitializeMainUser(rivalInfo *vo.InitializeRivalInfoVo) error {
	if rivalInfo.ImportStrategy == "" {
		return eris.New("init main user: import strategy cannot be empty")
	}
	if rivalInfo.Locale != nil && *rivalInfo.Locale != "" {
		conf, err := config.ReadConfig()
		if err != nil {
			return eris.Wrap(err, "read config")
		}
		conf.Locale = *rivalInfo.Locale
		if err := conf.WriteConfig(); err != nil {
			return eris.Wrap(err, "write config")
		}
	}
	mainUserCount, err := selectRivalInfoCount(s.db, &vo.RivalInfoVo{MainUser: true})
	if err != nil {
		return eris.Wrap(err, "select count from rival_info")
	}
	if mainUserCount > 0 {
		return eris.New("cannot have two main user, what are you doing?")
	}
	if len(rivalInfo.BMSDirectories) > 0 && rivalInfo.ImportStrategy == "LR2" {
		if err := saveSongDirectories(s.db, rivalInfo.BMSDirectories); err != nil {
			return eris.Wrap(err, "save directories")
		}
	}
	insertRivalInfo := rivalInfo.Into()
	insertRivalInfo.MainUser = true
	insertRivalInfo.ID = 1
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
	if insertRivalInfo.ScoreDataPath == nil || *insertRivalInfo.ScoreDataPath == "" {
		return eris.New("score.db path cannot be empty")
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return addRivalInfo(tx, insertRivalInfo)
	}); err != nil {
		return eris.Wrap(err, "transaction")
	}

	if insertRivalInfo.Type != entity.RIVAL_TYPE_LR2 {
		if err := s.monitorService.SetScoreLogFilePath(*insertRivalInfo.ScoreLogPath); err != nil {
			log.Errorf("monitor: cannot setup monitor service: %s", err)
		}
	}
	return nil
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
		return nil, eris.Wrap(err, "choose directory")
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
		return eris.Errorf("rival name cannot be empty")
	}
	// No, you can never add a main user by using this inteface
	rivalInfo.MainUser = false
	return eris.Wrap(s.db.Transaction(func(tx *gorm.DB) error {
		return addRivalInfo(tx, rivalInfo.Entity())
	}), "transaction")
}

func (s *RivalInfoService) FindRivalInfoList(filter *vo.RivalInfoVo) ([]*dto.RivalInfoDto, int, error) {
	if data, n, err := findRivalInfoList(s.db, filter); err != nil {
		return nil, 0, eris.Wrap(err, "findRivalInfoList")
	} else {
		return data, n, nil
	}
}

func (s *RivalInfoService) FindRivalInfoByID(rivalID uint) (*entity.RivalInfo, error) {
	if data, err := findRivalInfoByID(s.db, rivalID); err != nil {
		return nil, eris.Wrap(err, "findRivalInfoByID")
	} else {
		return data, nil
	}
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
			if err := s.FullySyncRivalData(rivalInfo); err != nil {
				return eris.Wrap(err, "fully sync rival data")
			}
		}
		log.Debug("[RivalInfoService] dispatched into incremental reload")
		if err := s.IncrementalSyncRivalData(rivalInfo); err != nil {
			return eris.Wrap(err, "incrementally sync rival data")
		}
	}
	return nil
}

// Fully reload one rival's save files
func (s *RivalInfoService) FullySyncRivalData(rivalInfo *entity.RivalInfo) error {
	if rivalInfo == nil {
		return eris.Errorf("assert: SyncRivalData: rivalInfo == nil")
	}
	if rivalInfo.ID == 0 {
		return eris.Errorf("assert: SyncRivalData: rivalInfo.ID corrupted")
	}
	log.Debug("[Service] calling RivalInfoService.SyncRivalData")
	if rivalInfo.ScoreLogPath == nil || *rivalInfo.ScoreLogPath == "" {
		return eris.Errorf("cannot sync rival %s's score log: score log file path is empty", rivalInfo.Name)
	}
	return eris.Wrap(s.db.Transaction(func(tx *gorm.DB) error {
		if err := syncRivalData(tx, rivalInfo, &vo.RivalFileReloadInfoVo{
			SongData:     true,
			ScoreLog:     true,
			ScoreDataLog: true,
			ScoreData:    true,
		}); err != nil {
			return eris.Wrap(err, "syncRivalData")
		}
		if err := updateRivalPlayCount(tx, rivalInfo.ID); err != nil {
			return eris.Wrap(err, "updateRivalPlayCount")
		}
		return nil
	}), "transaction")
}

// Extension to SyncRivalData, which only reloads part of the scorelog.db and
// scoredatalog.db file (if provided, or main user only). While songdata.db won't be read
// More specifically, only reloads the log that is set after rival's last log
// that recorded in lampghost
//
// Requirements:
//  1. rivalInfo's id > 0
//  2. rivalInfo's scorelog path must not be empty
//
// NOTE: This function's motivation is giving user a sync way that won't reload
// the songdata.db file but only the play log files. Since songdata.db won't have
// much changes often
func (s *RivalInfoService) IncrementalSyncRivalData(rivalInfo *entity.RivalInfo) error {
	if rivalInfo == nil {
		return eris.Errorf("incrementalSyncRivalData: rivalInfo cannot be nil")
	}
	if rivalInfo.ID == 0 {
		return eris.Errorf("incrementalSyncRivalData: rivalInfo.ID should > 0")
	}
	if rivalInfo.ScoreLogPath == nil || *rivalInfo.ScoreLogPath == "" {
		return eris.Errorf("incrementalSyncRivalData: rivalInfo.ScoreLogPath cannot be empty")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := incrementalSyncRivalData(tx, rivalInfo); err != nil {
			return eris.Wrap(err, "incrementalSyncRivalData")
		}
		if err := updateRivalPlayCount(tx, rivalInfo.ID); err != nil {
			return eris.Wrap(err, "updateRivalPlayCount")
		}
		return nil
	})
}

func (s *RivalInfoService) DelRivalInfo(ID uint) error {
	mainUser, err := queryMainUser(s.db)
	if err != nil {
		return eris.Wrap(err, "queryMainUser")
	}
	if ID == mainUser.ID {
		return eris.Errorf("DelRivalInfo: cannot delete main user")
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		var candidate entity.RivalInfo
		if err := tx.First(&candidate, ID).Error; err != nil {
			return eris.Wrap(err, "query deleting rival")
		}
		// RivalInfo
		if err := tx.Delete(&entity.RivalInfo{}, candidate.ID).Error; err != nil {
			return eris.Wrap(err, "delete rival")
		}
		// RivalScoreLog
		if err := tx.Unscoped().Where("rival_id = ?", candidate.ID).Delete(&entity.RivalScoreLog{}).Error; err != nil {
			return eris.Wrap(err, "delete rival_score_log")
		}
		// RivalSongData
		if err := tx.Unscoped().Where("rival_id = ?", candidate.ID).Delete(&entity.RivalSongData{}).Error; err != nil {
			return eris.Wrap(err, "delete rival_score_data_log")
		}
		// RivalTag
		if err := tx.Unscoped().Where("rival_id = ?", candidate.ID).Delete(entity.RivalTag{}).Error; err != nil {
			return eris.Wrap(err, "delete rival_tag")
		}
		return nil
	}); err != nil {
		return eris.Wrap(err, "transaction")
	}
	return nil
}

func (s *RivalInfoService) QueryUserPlayCountInYear(ID uint, yearNum string) ([]int, error) {
	logs, _, err := findRivalScoreLogList(s.db, &vo.RivalScoreLogVo{
		RivalId:     ID,
		SpecifyYear: &yearNum,
	})
	if err != nil {
		return nil, eris.Wrap(err, "findRivalScoreLogList")
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
		return nil, eris.Wrap(err, "FindRivalInfoByID")
	}
	header, err := queryLevelLayeredDiffTableInfoById(s.db, headerID)
	if err != nil {
		return nil, eris.Wrap(err, "queryLevelLayeredDiffTableInfoById")
	}
	sha256MaxLamp, err := findRivalMaximumClearScoreLogSha256Map(s.db, &vo.RivalScoreLogVo{
		RivalId: rivalID,
	})
	if err != nil {
		return nil, eris.Wrap(err, "findRivalMaximumClearScoreLogSha256Map")
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
		return nil, 0, eris.Wrap(err, "query rival_score_log")
	}
	return years, len(years), nil
}

func (s *RivalInfoService) QueryMainUser() (*entity.RivalInfo, error) {
	if data, err := queryMainUser(s.db); err != nil {
		return nil, eris.Wrap(err, "queryMainUser")
	} else {
		return data, nil
	}
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
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		var prev entity.RivalInfo
		if err := tx.Debug().First(&prev, rivalInfo.ID).Error; err != nil {
			return eris.Wrap(err, "query rival_info")
		}

		if prev.MainUser && (rivalInfo.SongDataPath == nil || *rivalInfo.SongDataPath == "") {
			return eris.Errorf("updateRivalInfo: SongDataPath cannot be empty when updateing main user")
		}

		if !prev.MainUser && (rivalInfo.SongDataPath != nil && *rivalInfo.SongDataPath != "") {
			return eris.Errorf("updateRivalInfo: cannot provide songdata path for a non main user")
		}
		rivalInfo.Type = prev.Type

		reloadInfo := &vo.RivalFileReloadInfoVo{
			SongData:     false,
			ScoreLog:     false,
			ScoreDataLog: false,
			ScoreData:    false,
		}
		if rivalInfo.ScoreLogPath != nil && *rivalInfo.ScoreLogPath != *prev.ScoreLogPath {
			reloadInfo.ScoreLog = true
		}
		if rivalInfo.SongDataPath != nil && *rivalInfo.SongDataPath != *prev.SongDataPath {
			reloadInfo.SongData = true
		}
		if rivalInfo.ScoreDataLogPath != nil && (prev.ScoreDataLogPath == nil || *prev.ScoreDataLogPath != *rivalInfo.ScoreDataLogPath) {
			reloadInfo.ScoreDataLog = true
		}
		if rivalInfo.ScoreDataPath != nil && (prev.ScoreDataPath == nil || *prev.ScoreDataPath != *rivalInfo.ScoreDataPath) {
			reloadInfo.ScoreData = true
		}
		if prev.MainUser && reloadInfo.ScoreLog {
			s.monitorService.SetScoreLogFilePath(*rivalInfo.ScoreLogPath)
		}
		if err := updateRivalInfo(tx, reloadInfo, rivalInfo.Entity()); err != nil {
			return eris.Wrap(err, "updateRivalInfo")
		}
		return nil
	}); err != nil {
		return eris.Wrap(err, "transaction")
	}
	return nil
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
	if err := updateRivalReverseImportInfo(s.db, *updateParam.Entity()); err != nil {
		return eris.Wrap(err, "update rival info plain fields")
	}
	return nil
}

func addRivalInfo(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
	if err := tx.Create(rivalInfo).Error; err != nil {
		return eris.Wrap(err, "insert rival_info")
	}
	if err := syncRivalData(tx, rivalInfo, &vo.RivalFileReloadInfoVo{
		SongData:     true,
		ScoreLog:     true,
		ScoreDataLog: true,
		ScoreData:    true,
	}); err != nil {
		return eris.Wrap(err, "syncRivalData")
	}
	return eris.Wrap(updateRivalPlayCount(tx, rivalInfo.ID), "updateRivalPlayCount")
}

func findRivalInfoList(tx *gorm.DB, filter *vo.RivalInfoVo) ([]*dto.RivalInfoDto, int, error) {
	var out []*dto.RivalInfoDto
	fields := `
		rival_info.*,
		rival_tag.tag_name
	`
	moved := tx.Select(fields).Model(&entity.RivalInfo{}).Scopes(scopeRivalInfoFilter(filter))
	moved = moved.Joins(`left join rival_tag on rival_info.lock_tag_id = rival_tag.id`)

	if err := moved.Debug().Find(&out).Error; err != nil {
		return nil, 0, eris.Wrap(err, "query rival_info")
	}

	return out, len(out), nil
}

func findRivalInfoByID(tx *gorm.DB, ID uint) (out *entity.RivalInfo, err error) {
	if err := tx.First(&out, ID).Error; err != nil {
		return nil, eris.Wrap(err, "query rival_info")
	}
	return
}

func selectRivalInfoCount(tx *gorm.DB, filter *vo.RivalInfoVo) (int64, error) {
	var count int64
	if err := tx.Model(&entity.RivalInfo{}).Scopes(scopeRivalInfoFilter(filter)).Count(&count).Error; err != nil {
		return 0, eris.Wrap(err, "query rival_info")
	}
	return count, nil
}

// Fully reload one rival's saves files
//
// WARN: This function wouldn't update rival's play count, it's not this
// function's purpose
//
// For now, theses files would be reloaded if required to
//  1. scorelog.db
//  2. songdata.db
//  3. scoredatalog.db
//
// NOTE: The way we re-generate rival's tags is based on whether songdata.db
// file exists or not, see below implementation for details.
func syncRivalData(tx *gorm.DB, rivalInfo *entity.RivalInfo, reloadInfo *vo.RivalFileReloadInfoVo) (err error) {
	// Quick quit if nothing specified
	if !reloadInfo.SongData && !reloadInfo.ScoreLog && !reloadInfo.ScoreDataLog && !reloadInfo.ScoreData {
		log.Debugf("no file needs to be reload, skipping...")
		return nil
	}
	// (1) songdata.db
	var songHashCache *entity.SongHashCache
	if reloadInfo.SongData && rivalInfo.SongDataPath != nil && *rivalInfo.SongDataPath != "" {
		var rawSongData []*entity.SongData
		rawSongData, err = loadSongData(*rivalInfo.SongDataPath)
		if err != nil {
			return eris.Wrap(err, "load song data")
		}
		if err := syncSongData(tx, rawSongData, rivalInfo.ID); err != nil {
			return eris.Wrap(err, "sync song data")
		}
		songHashCache = generateSongHashCacheFromRawData(rawSongData)
		expireDefaultCache()
	} else {
		var err error
		songHashCache, err = queryDefaultSongHashCache(tx)
		if err != nil {
			return eris.Wrap(err, "query default song hash cache")
		}
	}
	// (2) scorelog.db
	if reloadInfo.ScoreLog && rivalInfo.ScoreLogPath != nil && *rivalInfo.ScoreLogPath != "" {
		scoreLogService := NewScoreLogService(tx, rivalInfo.Type, false, songHashCache)
		if err := scoreLogService.SyncScoreLog(rivalInfo.ID, *rivalInfo.ScoreLogPath); err != nil {
			return eris.Wrap(err, "sync score log")
		}
	}
	// (3) scoredatalog.db
	if reloadInfo.ScoreDataLog && rivalInfo.ScoreDataLogPath != nil && *rivalInfo.ScoreDataLogPath != "" {
		scoreDataLogService := NewScoreDataLogService(tx, rivalInfo.Type, false, songHashCache)
		if err := scoreDataLogService.SyncScoreDataLog(rivalInfo.ID, *rivalInfo.ScoreDataLogPath); err != nil {
			return eris.Wrap(err, "sync score data log")
		}
	}
	// (4) score.db
	if reloadInfo.ScoreData && rivalInfo.ScoreDataPath != nil && *rivalInfo.ScoreDataPath != "" {
		scoreDataService := NewScoreDataService(tx, rivalInfo.Type, false, songHashCache)
		if err := scoreDataService.SyncScore(rivalInfo.ID, *rivalInfo.ScoreDataPath); err != nil {
			return eris.Wrap(err, "sync score data")
		}
	}
	// (4) generate rival tags
	return eris.Wrap(syncRivalTag(tx, rivalInfo.ID), "syncRivalTag")
}

// Incrementally reload one rival's save files
//
// WARN: This function wouldn't update rival's play count, it's not
// this function's purpose
//
// For now, theses files would be reloaded
//  1. scorelog.db
//  2. scoredatalog.db (if provides)
//  3. score.db (if provides)
//
// And these tables' data would be updated
//  1. rival_score_log (incrementally added for beatoraja)
//  2. rival_score_data_log (incrementally added for beatoraja)
//  3. rival_score_data (incrementally updated for beatoraja)
//  4. rival_tag (keep the old data as much as possible)
//
// NOTE: For LR2 user, rival_score_log and rival_score_data_log would be regenerated because
// LR2's database doesn't provide record's time
//
// NOTE: This function would degenerate to fully reload(which would trigger a songdata.db reload) if:
//  1. rival_score_log is empty
//  2. rival_score_data_log is empty (This wouldn't be a problem since songdata.db could only be imported by main user currently)
func incrementalSyncRivalData(tx *gorm.DB, rivalInfo *entity.RivalInfo) error {
	songHashCache, err := queryDefaultSongHashCache(tx)
	if err != nil {
		return eris.Wrap(err, "queryDefaultSongHashCache")
	}

	scoreLogService := NewScoreLogService(tx, rivalInfo.Type, true, songHashCache)
	if err := scoreLogService.SyncScoreLog(rivalInfo.ID, *rivalInfo.ScoreLogPath); err != nil {
		return eris.Wrap(err, "sync score log")
	}

	if rivalInfo.ScoreDataLogPath != nil && *rivalInfo.ScoreDataLogPath != "" {
		scoreDataLogService := NewScoreDataLogService(tx, rivalInfo.Type, true, songHashCache)
		if err := scoreDataLogService.SyncScoreDataLog(rivalInfo.ID, *rivalInfo.ScoreDataLogPath); err != nil {
			return eris.Wrap(err, "sync score data log")
		}
	}

	if rivalInfo.ScoreDataPath != nil && *rivalInfo.ScoreDataPath != "" {
		scoreDataService := NewScoreDataService(tx, rivalInfo.Type, true, songHashCache)
		if err := scoreDataService.SyncScore(rivalInfo.ID, *rivalInfo.ScoreDataPath); err != nil {
			return eris.Wrap(err, "sync score data")
		}
	}

	return eris.Wrap(syncRivalTag(tx, rivalInfo.ID), "syncRivalTag")
}

// Read all contents from `songdata.db` file into memory
func loadSongData(songDataPath string) ([]*entity.SongData, error) {
	if songDataPath == "" {
		return nil, eris.Errorf("assert: songdata path cannot be empty")
	}
	log.Debugf("[RivalInfoService] Trying to read song data from %s", songDataPath)
	if err := database.VerifyLocalDatabaseFilePath(songDataPath); err != nil {
		return nil, eris.Wrap(err, "verify local database file")
	}
	dsn := songDataPath + READONLY_PARAMETER
	songDataDB, err := gorm.Open(sqlite.Open(dsn))
	if err != nil {
		return nil, eris.Wrap(err, "open songdata.db")
	}
	songDataService := NewSongDataService(songDataDB)
	rawSongData, n, err := songDataService.FindSongDataList()
	if err != nil {
		return nil, eris.Wrap(err, "FindSongDataList")
	}
	log.Infof("[RivalInfoService] Read %d song data from %s", n, songDataPath)
	return rawSongData, nil
}

// Fully delete all content from rival_score_log and rebuild them by rawScoreLog
func syncSongData(tx *gorm.DB, rawSongData []*entity.SongData, rivalID uint) error {
	rivalSongData := Map(rawSongData, func(rawData *entity.SongData, _ int) *entity.RivalSongData {
		ret := entity.FromRawSongDataToRivalSongData(rawData)
		ret.RivalId = rivalID
		return &ret
	})

	if err := tx.Unscoped().Where("rival_id = ?", rivalID).Delete(&entity.RivalSongData{}).Error; err != nil {
		return eris.Wrap(err, "delete rival_song_data")
	}

	return eris.Wrap(tx.CreateInBatches(&rivalSongData, DEFAULT_BATCH_SIZE).Error, "insert rival_song_data")
}

func queryMainUser(tx *gorm.DB) (out *entity.RivalInfo, err error) {
	err = eris.Wrap(tx.Where(&entity.RivalInfo{MainUser: true}).First(&out).Error, "query rival_info")
	return
}

// Update one rival info
//
// Special Requirements:
//  1. when updating main user, songdata.db file path cannot be empty
func updateRivalInfo(tx *gorm.DB, reloadInfo *vo.RivalFileReloadInfoVo, rivalInfo *entity.RivalInfo) error {
	if err := syncRivalData(tx, rivalInfo, reloadInfo); err != nil {
		return eris.Wrap(err, "syncRivalData")
	}

	// Since we have reload scorelog.db, we need also update play count here
	if reloadInfo.ScoreLog {
		pc, err := selectRivalScoreLogCount(tx, &vo.RivalScoreLogVo{RivalId: rivalInfo.ID})
		if err != nil {
			return eris.Wrap(err, "selectRivalScoreLogCount")
		}
		rivalInfo.PlayCount = int(pc)
	}

	return eris.Wrap(tx.Updates(rivalInfo).Error, "update rival_info")
}

func updateRivalReverseImportInfo(tx *gorm.DB, rivalInfo entity.RivalInfo) error {
	updates := make(map[string]any)
	updates["ReverseImport"] = rivalInfo.ReverseImport
	updates["LockTagID"] = rivalInfo.LockTagID
	return eris.Wrap(tx.Model(&rivalInfo).Updates(updates).Error, "update rival_info")
}

// Simple helper function for updating one rival's play count field
func updateRivalPlayCount(tx *gorm.DB, rivalID uint) error {
	if rivalID == 0 {
		return eris.Errorf("updateRivalPlayCount: rivalID cannot be 0")
	}

	pc, err := selectRivalScoreLogCount(tx, &vo.RivalScoreLogVo{RivalId: rivalID})
	if err != nil {
		return err
	}
	return eris.Wrap(tx.Model(&entity.RivalInfo{Model: gorm.Model{ID: rivalID}}).Update("play_count", pc).Error, "update rival_info")
}

// Common query scope for vo.RivalInfoVo
func scopeRivalInfoFilter(filter *vo.RivalInfoVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Add extra filter here
		if filter.IgnoreMainUser {
			moved = moved.Where("rival_info.ID != 1")
		}
		return moved
	}
}
