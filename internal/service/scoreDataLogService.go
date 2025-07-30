package service

import (
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type ScoreDataLogService struct {
	source string // constants: "beatoraja" | "LR2"
}

func NewScoreDataLogService(source string) *ScoreDataLogService {
	return &ScoreDataLogService{
		source: source,
	}
}

// Load one scoredatalog.db/user.db as RivalScoreDataLog into memory
func (s *ScoreDataLogService) LoadScoreDataLog(rivalID uint, songHashCache *entity.SongHashCache, scoreDataLogPath string, maximumTimestamp *int64) ([]*entity.RivalScoreDataLog, int, error) {
	switch s.source {
	case entity.RIVAL_TYPE_BEATORAJA:
		rawLogs, err := loadScoreDataLog(scoreDataLogPath, maximumTimestamp)
		if err != nil {
			return nil, 0, err
		}
		rivalScoreDataLog := make([]*entity.RivalScoreDataLog, len(rawLogs))
		for i, rawLog := range rawLogs {
			dataLog := entity.FromRawScoreDataLogToRivalScoreDataLog(rawLog)
			dataLog.RivalId = rivalID
			rivalScoreDataLog[i] = &dataLog
		}
		return rivalScoreDataLog, len(rivalScoreDataLog), nil
	case entity.RIVAL_TYPE_LR2:
		rawLogs, err := loadLR2Log(scoreDataLogPath)
		if err != nil {
			return nil, 0, err
		}
		// NOTE: See LoadScoreLog for details
		rivalScoreDataLog := make([]*entity.RivalScoreDataLog, len(rawLogs))
		for i, rawLog := range rawLogs {
			dataLog := entity.FromRawLR2LogToRivalScoreDataLog(rawLog)
			if sha256, ok := songHashCache.GetSHA256(rawLog.MD5); ok {
				dataLog.Sha256 = sha256
			}
			dataLog.RivalId = rivalID
			rivalScoreDataLog[i] = &dataLog
		}
		return rivalScoreDataLog, len(rivalScoreDataLog), nil
	default:
		panic("unexpected source: " + s.source)
	}
}

// Read all contents from `scoredatalog.db` (beatoraja) file into memory
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
	rawScoreDataLog, n, err := findBeatorajaScoreDataLogList(scoreDataLogDB, maximumTimestamp)
	if err != nil {
		return nil, eris.Wrap(err, "load: query from scoredatalog.db failed")
	}
	log.Debugf("[RivalInfoService] Read %d logs from %s", n, scoreDataLogPath)
	return rawScoreDataLog, nil
}

func findBeatorajaScoreDataLogList(tx *gorm.DB, maximumTimestamp *int64) ([]*entity.ScoreDataLog, int, error) {
	var data []*entity.ScoreDataLog
	moved := tx
	if maximumTimestamp != nil {
		moved = moved.Where("date > ?", maximumTimestamp)
	}
	if err := moved.Find(&data).Error; err != nil {
		return nil, 0, err
	}
	return data, len(data), nil
}
