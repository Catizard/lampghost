package service

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type ScoreLogService struct {
	source string // constants: "beatoraja" | "LR2"
}

func NewScoreLogService(source string) *ScoreLogService {
	return &ScoreLogService{
		source: source,
	}
}

// Load one scorelog.db/user.db as RivalScoreLog into memory
func (s *ScoreLogService) LoadScoreLog(rivalID uint, songHashCache *entity.SongHashCache, scoreLogPath string, maximumTimestamp *int64) ([]*entity.RivalScoreLog, int, error) {
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
		for i, rawLog := range rawLogs {
			rivalLog := entity.FromRawLR2LogToRivalScoreLog(rawLog)
			rivalLog.RivalId = rivalID
			if sha256, ok := songHashCache.GetSHA256(rawLog.MD5); ok {
				rivalLog.Sha256 = sha256
			}
			rivalScoreLog[i] = &rivalLog
		}
		return rivalScoreLog, len(rivalScoreLog), nil
	default:
		panic("unexpected source: " + s.source)
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
