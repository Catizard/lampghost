package service

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"gorm.io/gorm"
)

type RivalScoreLogService struct {
	db *gorm.DB
}

func NewRivalScoreLogService(db *gorm.DB) *RivalScoreLogService {
	return &RivalScoreLogService{
		db: db,
	}
}

func (s *RivalScoreLogService) QueryRivalScoreLogPageList(filter *vo.RivalScoreLogVo) ([]*dto.RivalScoreLogDto, int, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var count int64
	if err := tx.Model(&entity.RivalScoreLog{}).Count(&count).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}
	scoreLogs, n, err := findRivalScoreLogList(tx, filter)
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, 0, err
	}
	return scoreLogs, n, nil
}

func findRivalScoreLogList(tx *gorm.DB, filter *vo.RivalScoreLogVo) ([]*dto.RivalScoreLogDto, int, error) {
	fields := `
		rival_score_log.*,
		sd.title as title,
		sd.md5 as md5,
		sd.ID as rival_song_data_id
	`
	partial := tx.Model(&entity.RivalScoreLog{}).Order("rival_score_log.record_time desc").Select(fields)
	var out []*dto.RivalScoreLogDto
	// TODO: left join on rival_song_data is the bottleneck, how to replace it?
	partial = partial.Debug().Joins("left join (select * from rival_song_data group by sha256) as sd on rival_score_log.sha256 = sd.sha256").Scopes(
		scopeRivalScoreLogFilter(filter),
		pagination(filter.Pagination),
	)
	if filter.SongNameLike != nil && *filter.SongNameLike != "" {
		partial = partial.Where("sd.title like ('%' || ? || '%')", filter.SongNameLike)
	}
	if err := partial.Find(&out).Error; err != nil {
		return nil, 0, err
	}
	// pagination
	if filter != nil && filter.Pagination != nil {
		count, err := selectRivalScoreLogCount(tx, filter)
		if err != nil {
			return nil, 0, err
		}
		filter.Pagination.PageCount = calcPageCount(count, filter.Pagination.PageSize)
	}
	return out, len(out), nil
}

// Query the last played score log
func findLastRivalScoreLogList(tx *gorm.DB, filter *vo.RivalScoreLogVo) (*entity.RivalScoreLog, error) {
	ret := entity.RivalScoreLog{}
	err := tx.Model(&entity.RivalScoreLog{}).Scopes(scopeRivalScoreLogFilter(filter)).Order("record_time desc").Limit(1).Find(&ret).Error
	return &ret, err
}

// Extend function to findRivalScoreLogList
//
// Returns sha256 grouped array
func findRivalScoreLogSha256Map(tx *gorm.DB, filter *vo.RivalScoreLogVo) (map[string][]*dto.RivalScoreLogDto, error) {
	scoreLogs, _, err := findRivalScoreLogList(tx, filter)
	if err != nil {
		return nil, err
	}
	sha256ScoreLogsMap := make(map[string][]*dto.RivalScoreLogDto)
	for _, scoreLog := range scoreLogs {
		if _, ok := sha256ScoreLogsMap[scoreLog.Sha256]; !ok {
			sha256ScoreLogsMap[scoreLog.Sha256] = make([]*dto.RivalScoreLogDto, 0)
		}
		sha256ScoreLogsMap[scoreLog.Sha256] = append(sha256ScoreLogsMap[scoreLog.Sha256], scoreLog)
	}
	return sha256ScoreLogsMap, nil
}

func selectRivalScoreLogCount(tx *gorm.DB, filter *vo.RivalScoreLogVo) (int64, error) {
	if filter == nil {
		var count int64
		if err := tx.Debug().Model(&entity.RivalScoreLog{}).Count(&count).Error; err != nil {
			return 0, err
		}
		return count, nil
	}
	var count int64
	if err := tx.Model(&entity.RivalScoreLog{}).Debug().Where(filter.Entity()).Scopes(
		scopeRivalScoreLogFilter(filter),
	).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Specialized scope for vo.RivalScoreLogVo
func scopeRivalScoreLogFilter(filter *vo.RivalScoreLogVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Extra filters
		if filter.OnlyCourseLogs {
			moved = moved.Where("length(rival_score_log.sha256) > 64")
		}
		if !filter.StartRecordTime.IsZero() {
			moved = moved.Where("rival_score_log.record_time >= ?", filter.StartRecordTime)
		}
		if !filter.EndRecordTime.IsZero() {
			moved = moved.Where("rival_score_log.record_time <= ?", filter.EndRecordTime)
		}
		if filter.MinimumClear != nil {
			moved = moved.Where("rival_score_log.clear >= ?", filter.MinimumClear)
		}
		if filter.SpecifyYear != nil {
			moved = moved.Where("STRFTIME('%Y', `rival_score_log`.`record_time`) = ?", filter.SpecifyYear)
		}
		return moved
	}
}
