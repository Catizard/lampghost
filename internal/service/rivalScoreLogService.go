package service

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/rotisserie/eris"
	. "github.com/samber/lo"
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
	tableTags, _, err := queryDiffTableTag(tx, &vo.DiffTableDataVo{
		Md5s: Map(scoreLogs, func(log *dto.RivalScoreLogDto, _ int) string {
			return log.Md5
		}),
	})
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}
	ForEach(scoreLogs, func(log *dto.RivalScoreLogDto, _ int) {
		log.TableTags = make([]*dto.DiffTableTagDto, 0)
		ForEach(tableTags, func(tag *dto.DiffTableTagDto, _ int) {
			if tag.Md5 == log.Md5 {
				log.TableTags = append(log.TableTags, tag)
			}
		})
	})
	if err := tx.Commit().Error; err != nil {
		return nil, 0, err
	}
	return scoreLogs, n, nil
}

// Return play logs in one specified day, which satisfy:
//  1. It's strict lower than filter.EndRecordTime
//  2. It's the maximum possible one
//
// The frontend should call this function as follow:
//  1. Pass now() + 1 day
//  2. If return value is not empty, pick record time from anyone and set it to next parameter
//  3. If return value is empty, mark no more values and stop querying
//  4. Or repeat querying the next stocks
//
// TODO: This function actually drops course challenge logs, it's a lilltle hard to write it
// correctly so I decided to leave this feature as a todo.
//
// NOTE: Result is unique, only the maximum clear one would be reserved. And it's ordered,
// the higher clear one would be placed before the lower one.
func (s *RivalScoreLogService) QueryPrevDayScoreLogList(filter *vo.RivalScoreLogVo) ([]*dto.RivalScoreLogDto, int, error) {
	rawLogs, _, err := queryPrevDayScoreLogList(s.db, filter)
	if err != nil {
		return nil, 0, err
	}
	logs := make([]*dto.RivalScoreLogDto, 0)
	for i := 0; i < len(rawLogs); i++ {
		j := i
		for j+1 < len(rawLogs) && rawLogs[j+1].ID == rawLogs[i].ID {
			j++
		}
		log := rawLogs[i]
		log.TableTags = make([]*dto.DiffTableTagDto, 0)
		if log.TableName != "" {
			for k := i; k <= j; k++ {
				log.TableTags = append(log.TableTags, dto.PushDownTag(rawLogs[k]))
			}
		}
		logs = append(logs, log)
		i = j
	}
	return logs, len(logs), nil
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

// Like findRivalScoreLogList, but only reserve the maximum clear one
func findRivalMaximumClearScoreLogList(tx *gorm.DB, filter *vo.RivalScoreLogVo) ([]*dto.RivalScoreLogDto, int, error) {
	fields := `
		rival_score_log.*,
		max(rival_score_log.clear) as clear
	`
	partial := tx.Model(&entity.RivalScoreLog{}).Order("record_time desc").Select(fields)
	var out []*dto.RivalScoreLogDto
	partial = partial.Scopes(scopeRivalScoreLogFilter(filter))
	if filter != nil {
		partial = partial.Scopes(scopeInSha256s(filter.Sha256s))
	}
	// NOTE: without this statement, this function has strange behaviour
	partial = partial.Group("sha256")
	if err := partial.Debug().Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, len(out), nil
}

// Extend function to findRivalMaximumClearScoreLogList
//
// Returns sha256 grouped array
func findRivalMaximumClearScoreLogSha256Map(tx *gorm.DB, filter *vo.RivalScoreLogVo) (map[string][]*dto.RivalScoreLogDto, error) {
	scorelogs, _, err := findRivalMaximumClearScoreLogList(tx, filter)
	if err != nil {
		return nil, err
	}
	return groupingScoreLogsBySha256(scorelogs), err
}

// Query the last played score log
func findLastRivalScoreLog(tx *gorm.DB, filter *vo.RivalScoreLogVo) (*entity.RivalScoreLog, error) {
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
	return groupingScoreLogsBySha256(scoreLogs), nil
}

// Group score logs to a map, which key is sha256
func groupingScoreLogsBySha256(scoreLogs []*dto.RivalScoreLogDto) map[string][]*dto.RivalScoreLogDto {
	sha256ScoreLogsMap := make(map[string][]*dto.RivalScoreLogDto)
	for _, scoreLog := range scoreLogs {
		if _, ok := sha256ScoreLogsMap[scoreLog.Sha256]; !ok {
			sha256ScoreLogsMap[scoreLog.Sha256] = make([]*dto.RivalScoreLogDto, 0)
		}
		sha256ScoreLogsMap[scoreLog.Sha256] = append(sha256ScoreLogsMap[scoreLog.Sha256], scoreLog)
	}
	return sha256ScoreLogsMap
}

// NOTE: selectRivalScoreLogCount's filter statment should always be equal to selectRivalScoreLog
func selectRivalScoreLogCount(tx *gorm.DB, filter *vo.RivalScoreLogVo) (int64, error) {
	if filter == nil {
		var count int64
		if err := tx.Debug().Model(&entity.RivalScoreLog{}).Count(&count).Error; err != nil {
			return 0, err
		}
		return count, nil
	}
	var count int64
	partial := tx.Model(&entity.RivalScoreLog{}).Joins("left join (select * from rival_song_data group by sha256) as sd on rival_score_log.sha256 = sd.sha256").Scopes(
		scopeRivalScoreLogFilter(filter),
	)
	if filter.SongNameLike != nil && *filter.SongNameLike != "" {
		partial = partial.Where("sd.title like ('%' || ? || '%')", filter.SongNameLike)
	}
	if err := partial.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Warning: If one play log's song is related to multiple tables, the result would be duplicated.
// It's the caller side's job to group the same play logs and construct the tags correctly.
func queryPrevDayScoreLogList(tx *gorm.DB, filter *vo.RivalScoreLogVo) ([]*dto.RivalScoreLogDto, int, error) {
	fields := `
		rival_score_log.*,
		max(rival_score_log.clear),
		sd.title as title,
		sd.md5 as md5,
		sd.ID as rival_song_data_id,
		dt.table_name,
		dt.table_level,
		dt.table_symbol,
		dt.table_tag_color,
		dt.table_tag_text_color
	`
	// NOTE: Some filter statements should be applied to both subquery and the outer one
	partial := tx.Model(&entity.RivalScoreLog{}).Select(fields).Joins("left join (select ID, title, sha256, md5 from rival_song_data group by sha256) as sd on rival_score_log.sha256 = sd.sha256")
	// TODO: pagination rival_score_log has another implement way of tags, should be unified in the future
	partial = partial.Joins(`
		left join (
			select dd.md5, dh.name as table_name, dh.tag_color as table_tag_color, dh.tag_text_color as table_tag_text_color, dh.symbol as table_symbol, dd."level" as table_level
			from difftable_data dd
			left join difftable_header dh on dd.header_id = dh.id
			where dh.no_tag_build = 0
			group by dd.md5
		) as dt on sd.md5 = dt.md5
	`)
	// TODO: This sql is stricted with "clear >= 4", should we make this configurable?
	maximumDateQuery := tx.Model(&entity.RivalScoreLog{}).
		Select("date(max(rival_score_log.record_time))").
		Where("length(rival_score_log.sha256) = 64 and rival_score_log.clear >= 4 and date(record_time) < date(?) and rival_id = ?", filter.EndRecordTime, filter.RivalId)
	var out []*dto.RivalScoreLogDto
	partial = partial.Where("length(rival_score_log.sha256) = 64 and rival_score_log.clear >= 4 and rival_id = ?", filter.RivalId).Order("rival_score_log.clear desc")
	if err := partial.Debug().Where("date(rival_score_log.record_time) = (?)", maximumDateQuery).Group("rival_score_log.sha256").Find(&out).Error; err != nil {
		return nil, 0, eris.Wrap(err, "failed to query prev day score log")
	}
	return out, len(out), nil
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
		if filter.NoCourseLog {
			moved = moved.Where("length(rival_score_log.sha256) = 64")
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
