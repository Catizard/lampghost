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

type ScoreDataService struct {
	tx                  *gorm.DB              // lampghost database connection
	source              string                // constants: "beatoraja" | "LR2"
	incrementallyReload bool                  // when flagged, only append "new logs" if possible
	songHashCache       *entity.SongHashCache // converting md5 -> sha256 for LR2
}

func NewScoreDataService(tx *gorm.DB, source string, incrementallyReload bool, songHashCache *entity.SongHashCache) *ScoreDataService {
	return &ScoreDataService{
		tx:                  tx,
		source:              source,
		incrementallyReload: incrementallyReload,
		songHashCache:       songHashCache,
	}
}

// Sync one score.db/user.db into database
func (s *ScoreDataService) SyncScore(rivalID uint, scorePath string) error {
	var maximumTimestamp *int64 = nil
	if s.incrementallyReload {
		log.Debug("[ScoreDataService] call side requires an incrementally reload")
		switch s.source {
		case entity.RIVAL_TYPE_BEATORAJA:
			lastRivalScoreData, err := findLastRivalScoreData(s.tx, &vo.RivalScoreDataVo{RivalID: rivalID})
			if err != nil {
				return eris.Wrap(err, "find last rival score data")
			}
			if lastRivalScoreData.RivalID == 0 {
				s.incrementallyReload = false
				log.Debug("[ScoreDataService] degenerate to fully reload because no log is found")
			} else {
				timestamp := lastRivalScoreData.RecordTime.Unix()
				maximumTimestamp = &timestamp
				log.Debugf("[ScoreDataService] last record time is %v", lastRivalScoreData.RecordTime)
			}
		case entity.RIVAL_TYPE_LR2:
			s.incrementallyReload = false
		default:
			return eris.Errorf("unexpected rival type: %s", s.source)
		}
	}
	rawScoreData, _, err := s.LoadScoreData(rivalID, scorePath, maximumTimestamp)
	if err != nil {
		return eris.Wrap(err, "load score.db")
	}
	if s.incrementallyReload {
		if err := updateScoreData(s.tx, rawScoreData); err != nil {
			return eris.Wrap(err, "updateScoreData")
		}
	} else {
		if err := syncScoreData(s.tx, rawScoreData, rivalID); err != nil {
			return eris.Wrap(err, "syncScoreData")
		}
	}
	return nil
}

// Load one score.db/user.db as RivalScoreData into memory
func (s *ScoreDataService) LoadScoreData(rivalID uint, scoreDataPath string, maximumTimestamp *int64) ([]*entity.RivalScoreData, int, error) {
	switch s.source {
	case entity.RIVAL_TYPE_BEATORAJA:
		rawScoreData, err := loadScoreData(scoreDataPath, maximumTimestamp)
		if err != nil {
			return nil, 0, eris.Wrap(err, "loadScoreData")
		}
		rivalScoreData := make([]*entity.RivalScoreData, len(rawScoreData))
		for i, rawData := range rawScoreData {
			scoreData := entity.FromRawScoreDataToRivalScoreData(rawData)
			scoreData.RivalID = rivalID
			rivalScoreData[i] = &scoreData
		}
		return rivalScoreData, len(rivalScoreData), nil
	case entity.RIVAL_TYPE_LR2:
		rawLogs, err := loadLR2Log(scoreDataPath)
		if err != nil {
			return nil, 0, err
		}
		rivalScoreData := make([]*entity.RivalScoreData, len(rawLogs))
		for i, rawLog := range rawLogs {
			scoreData := entity.FromRawLR2LogToRivalScoreData(rawLog)
			if sha256, ok := s.songHashCache.GetSHA256(rawLog.MD5); ok {
				scoreData.Sha256 = sha256
			}
			scoreData.RivalID = rivalID
			rivalScoreData[i] = &scoreData
		}
		return rivalScoreData, len(rivalScoreData), nil
	default:
		return nil, 0, eris.Errorf("unexpected source: %s", s.source)
	}
}

// Read all contents from 'score.db' (beatoraja) file into memory
func loadScoreData(scoreDataPath string, maximumTimestamp *int64) ([]*entity.ScoreData, error) {
	if scoreDataPath == "" {
		return nil, eris.New("load: scoredata.db file path cannot be empty")
	}
	log.Debugf("[ScoreDataService] Trying to read log from %s", scoreDataPath)
	if err := database.VerifyLocalDatabaseFilePath(scoreDataPath); err != nil {
		return nil, err
	}
	dsn := scoreDataPath + READONLY_PARAMETER
	scoreDataDB, err := gorm.Open(sqlite.Open(dsn))
	if err != nil {
		return nil, eris.Wrap(err, "failed to open score.db")
	}
	rawScoreData, n, err := findBeatorajaScoreDataList(scoreDataDB, maximumTimestamp)
	if err != nil {
		return nil, eris.Wrap(err, "load: query from score.db failed")
	}
	log.Debugf("[ScoreDataService] Read %d logs from %s", n, scoreDataPath)
	return rawScoreData, nil
}

func findBeatorajaScoreDataList(tx *gorm.DB, maximumTimestamp *int64) ([]*entity.ScoreData, int, error) {
	var data []*entity.ScoreData
	moved := tx
	if maximumTimestamp != nil {
		moved = moved.Where("date > ?", maximumTimestamp)
	}
	if err := moved.Find(&data).Error; err != nil {
		return nil, 0, err
	}
	return data, len(data), nil
}
