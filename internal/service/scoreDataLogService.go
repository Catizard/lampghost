package service

import (
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type ScoreDataLogService struct {
	tx                  *gorm.DB              // lampghost database connection
	source              string                // constants: "beatoraja" | "LR2"
	incrementallyReload bool                  // when flagged, only append "new logs" if possible
	songHashCache       *entity.SongHashCache // converting md5 -> sha256 for LR2
}

func NewScoreDataLogService(tx *gorm.DB, source string, incrementallyReload bool, songHashCache *entity.SongHashCache) *ScoreDataLogService {
	return &ScoreDataLogService{
		tx:                  tx,
		source:              source,
		incrementallyReload: incrementallyReload,
		songHashCache:       songHashCache,
	}
}

// Sync one scoredatalog.db/user.db into database
func (s *ScoreDataLogService) SyncScoreDataLog(rivalID uint, scoreDataLogPath string) error {
	var maximumTimestamp *int64 = nil
	if s.incrementallyReload {
		log.Debug("[ScoreDataLogService] call side requires an incrementally reload")
		// NOTE: Degenerate to fully reload if no log is present or rival type is LR2
		switch s.source {
		case entity.RIVAL_TYPE_BEATORAJA:
			lastRivalScoreDataLog, err := findLastRivalScoreDataLog(s.tx, &vo.RivalScoreDataLogVo{RivalId: rivalID})
			if err != nil {
				return eris.Wrap(err, "find last rival score log")
			}
			if lastRivalScoreDataLog.ID == 0 {
				s.incrementallyReload = false
				log.Debug("[ScoreDataLogService] degenerate to fully reload because no log is found")
			} else {
				timestamp := lastRivalScoreDataLog.RecordTime.Unix()
				maximumTimestamp = &timestamp
				log.Debugf("[ScoreDataLogService] last record time is %v", lastRivalScoreDataLog.RecordTime)
			}
		case entity.RIVAL_TYPE_LR2:
			s.incrementallyReload = false
		default:
			return eris.Errorf("unexpected rival type: %s", s.source)
		}
	} else {
		log.Debug("[ScoreLogService] call side requires a fully reload")
	}
	rawScoreDataLog, _, err := s.LoadScoreDataLog(rivalID, scoreDataLogPath, maximumTimestamp)
	if err != nil {
		return eris.Wrap(err, "load scoredatalog.db")
	}
	if s.incrementallyReload {
		if err := appendScoreDataLog(s.tx, rawScoreDataLog); err != nil {
			return eris.Wrap(err, "append rival score data log")
		}
	} else {
		if err := syncScoreDataLog(s.tx, rawScoreDataLog, rivalID); err != nil {
			return eris.Wrap(err, "rebuild rival score data log")
		}
	}
	return nil
}

// Load one scoredatalog.db/user.db as RivalScoreDataLog into memory
func (s *ScoreDataLogService) LoadScoreDataLog(rivalID uint, scoreDataLogPath string, maximumTimestamp *int64) ([]*entity.RivalScoreDataLog, int, error) {
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
			if sha256, ok := s.songHashCache.GetSHA256(rawLog.MD5); ok {
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
