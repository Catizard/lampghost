package service

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type ScoreLogService struct {
	tx                  *gorm.DB              // lampghost database connection
	source              string                // constants: "beatoraja" | "LR2"
	incrementallyReload bool                  // when flagged, only append "new logs" if possible
	songHashCache       *entity.SongHashCache // converting md5 -> sha256 for LR2
}

func NewScoreLogService(tx *gorm.DB, source string, incrementallyReload bool, songHashCache *entity.SongHashCache) *ScoreLogService {
	return &ScoreLogService{
		tx:                  tx,
		source:              source,
		incrementallyReload: incrementallyReload,
		songHashCache:       songHashCache,
	}
}

// SyncScoreLog Sync one scorelog.db/user.db into database
func (s *ScoreLogService) SyncScoreLog(rivalID uint, scoreLogPath string) error {
	var maximumTimestamp *int64 = nil
	if s.incrementallyReload {
		log.Debug("[ScoreLogService] call side requires an incrementally reload")
		// NOTE: Degenerate to fully reload if no log is present or rival type is LR2
		switch s.source {
		case entity.RIVAL_TYPE_BEATORAJA:
			lastRivalScoreLog, err := findLastRivalScoreLog(s.tx, &vo.RivalScoreLogVo{RivalId: rivalID})
			if err != nil {
				return eris.Wrap(err, "find last rival score log")
			}
			if lastRivalScoreLog.ID == 0 {
				s.incrementallyReload = false
				log.Debug("[ScoreLogService] degenerate to fully reload because no log is found")
			} else {
				timestamp := lastRivalScoreLog.RecordTime.Unix()
				maximumTimestamp = &timestamp
				log.Debugf("[ScoreLogService] last record time is %v", lastRivalScoreLog.RecordTime)
			}
		case entity.RIVAL_TYPE_LR2:
			s.incrementallyReload = false
		default:
			return eris.Errorf("unexpected rival type: %s", s.source)
		}
	} else {
		log.Debug("[ScoreLogService] call side requires a fully reload")
	}
	rawScoreLog, _, err := s.LoadScoreLog(rivalID, scoreLogPath, maximumTimestamp)
	if err != nil {
		return eris.Wrap(err, "load scorelog.db")
	}
	if s.incrementallyReload {
		if err := appendScoreLog(s.tx, rawScoreLog); err != nil {
			return eris.Wrap(err, "append rival score log")
		}
	} else {
		if err := syncScoreLog(s.tx, rawScoreLog, rivalID); err != nil {
			return eris.Wrap(err, "rebuild rival score log")
		}
	}
	return nil
}

// LoadScoreLog Load one scorelog.db/user.db as RivalScoreLog into memory
func (s *ScoreLogService) LoadScoreLog(rivalID uint, scoreLogPath string, maximumTimestamp *int64) ([]*entity.RivalScoreLog, int, error) {
	switch s.source {
	case entity.RIVAL_TYPE_BEATORAJA:
		rawLogs, err := loadBeatorajaScoreLog(scoreLogPath, maximumTimestamp)
		if err != nil {
			return nil, 0, err
		}
		rivalScoreLog := make([]*entity.RivalScoreLog, len(rawLogs))
		for i, rawLog := range rawLogs {
			rivalLog := entity.FromRawScoreLogToRivalScoreLog(rawLog)
			rivalLog.RivalId = rivalID
			rivalScoreLog[i] = &rivalLog
		}
		return rivalScoreLog, len(rivalScoreLog), nil
	case entity.RIVAL_TYPE_LR2:
		// NOTE: LR2's table doesn't record when did the record set, therefore, it's impossible
		// to load the LR2's database incrementally. This is the reason why 'maximumTimestamp'
		// is ignored here
		rawLogs, err := loadLR2Log(scoreLogPath)
		if err != nil {
			return nil, 0, err
		}
		// NOTE: In order to make the recent activity component works for LR2 users, we implement
		// a diff algorithm here: each time we load the lr2's database, we calculate the difference
		// between current and the state in database based on lamp. If one song's clear lamp has
		// changed, we insert it as a rival_score_log row into database. This generates the data
		// that recent activity needs
		rivalScoreData, _, err := findRivalScoreDataList(s.tx, &vo.RivalScoreDataVo{
			RivalID: rivalID,
		})
		if err != nil {
			return nil, 0, err
		}
		oldLamps := make(map[string]int)
		for _, oldScoreData := range rivalScoreData {
			if old, ok := oldLamps[oldScoreData.Md5]; ok {
				// It's possible because the primary key is (md5, mode)
				oldLamps[oldScoreData.Md5] = max(old, oldLamps[oldScoreData.Md5])
			} else {
				oldLamps[oldScoreData.Md5] = old
			}
		}
		diffLogs := make([]*entity.LR2Log, 0)
		for _, rawLog := range rawLogs {
			if oldLamp, ok := oldLamps[rawLog.MD5]; ok {
				if oldLamp < rawLog.Clear {
					diffLogs = append(diffLogs, rawLog)
				}
			} else {
				diffLogs = append(diffLogs, rawLog)
			}
		}
		rivalScoreLog := make([]*entity.RivalScoreLog, len(rawLogs))
		// NOTE: Unlike beatoraja' scorelog, LR2's log is using MD5 as the hash field. There're
		// two ways to handle this:
		//  1. Add 'md5' field into 'rival_score_log' table, and rewrites all the sql that uses
		//     the sha256 before
		//  2. Convert 'md5' into 'sha256' before inserting into database, fully mimick the
		//     beatoraja's scorelog structure
		//
		// It's obvious that the first appoarch is nearly unaccpectable, it introduces too
		// much diverges and the implementation could be complicated. The second appoarch is less
		// damaging but it requires:
		//  1. song hash cache for converting md5 to sha256
		//  2. scan & repair the missing logs every time after updating song hash cache
		//
		// We still have to add a 'md5' field to 'rival_score_log' for repairing the hash lost logs,
		// and this issue is completely unseeable for users, which might be annoying
		for i, rawLog := range diffLogs {
			rivalLog := entity.FromRawLR2LogToRivalScoreLog(rawLog)
			rivalLog.RivalId = rivalID
			if sha256, ok := s.songHashCache.GetSHA256(rawLog.MD5); ok {
				rivalLog.Sha256 = sha256
			}
			rivalScoreLog[i] = &rivalLog
		}
		return rivalScoreLog, len(rivalScoreLog), nil
	default:
		return nil, 0, eris.Errorf("unexpected source: %s", s.source)
	}
}

func findBeatorajaScoreLogList(tx *gorm.DB, maximumTimestamp *int64) ([]*entity.ScoreLog, int, error) {
	var logs []*entity.ScoreLog
	moved := tx
	if maximumTimestamp != nil {
		moved = moved.Where("date > ?", maximumTimestamp)
	}
	if err := moved.Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, len(logs), nil
}

func findLR2LogList(tx *gorm.DB) ([]*entity.LR2Log, int, error) {
	var logs []*entity.LR2Log
	moved := tx
	if err := moved.Select("*, oid as row_id").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, len(logs), nil
}

// Read one `scorelog.db` (beatoraja) file into memory
func loadBeatorajaScoreLog(scoreLogPath string, maximumTimestamp *int64) ([]*entity.ScoreLog, error) {
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
	rawScoreLog, n, err := findBeatorajaScoreLogList(scoreLogDB, maximumTimestamp)
	if err != nil {
		return nil, err
	}
	log.Debugf("[RivalInfoService] Read %d logs from %s", n, scoreLogPath)
	return rawScoreLog, nil
}

// Read one `user.db` (LR2) file into memory
func loadLR2Log(userDBPath string) ([]*entity.LR2Log, error) {
	if userDBPath == "" {
		return nil, eris.New("assert: user.db path cannot be empty")
	}
	log.Debugf("[RivalInfoService] Trying to read user.db from %s", userDBPath)
	if err := database.VerifyLocalDatabaseFilePath(userDBPath); err != nil {
		return nil, eris.Wrap(err, "verify database path")
	}
	dsn := userDBPath + READONLY_PARAMETER
	log.Debugf("[RivalInfoService] full user.db path: %s", dsn)
	userDB, err := gorm.Open(sqlite.Open(dsn))
	if err != nil {
		return nil, eris.Wrap(err, "open user.db")
	}
	rawLR2Log, n, err := findLR2LogList(userDB)
	if err != nil {
		return nil, err
	}
	log.Debugf("[RivalInfoService] Read %d logs from %s", n, userDBPath)
	return rawLR2Log, nil
}
