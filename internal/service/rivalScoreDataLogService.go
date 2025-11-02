package service

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/random"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"gorm.io/gorm"

	. "github.com/samber/lo"
)

type RivalScoreDataLogService struct {
	db *gorm.DB
}

func NewRivalScoreDataLogService(db *gorm.DB) *RivalScoreDataLogService {
	return &RivalScoreDataLogService{
		db: db,
	}
}

func (s *RivalScoreDataLogService) QueryUserKeyCountInYear(param *vo.RivalScoreDataLogVo) (out []*dto.KeyCountDto, n int, err error) {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		out, n, err = queryUserKeyCountInYear(tx, param)
		return nil
	}); err != nil {
		return nil, 0, err
	}
	return out, n, nil
}

func (s *RivalScoreDataLogService) QueryRivalScoreDataLogPageList(filter *vo.RivalScoreDataLogVo) (out []*dto.RivalScoreDataLogDto, n int, err error) {
	err = s.db.Transaction(func(tx *gorm.DB) error {
		out, n, err = findRivalScoreDataLogList(tx, filter)
		if err != nil {
			return err
		}
		tableTags, _, err := queryDiffTableTag(tx, &vo.DiffTableDataVo{
			Md5s: FilterMap(out, func(log *dto.RivalScoreDataLogDto, _ int) (string, bool) {
				return log.Md5, log.Md5 != ""
			}),
		})
		if err != nil {
			return err
		}
		ForEach(out, func(log *dto.RivalScoreDataLogDto, _ int) {
			log.TableTags = make([]*dto.DiffTableTagDto, 0)
			ForEach(tableTags, func(tag *dto.DiffTableTagDto, _ int) {
				if log.Md5 != "" && tag.Md5 == log.Md5 {
					log.TableTags = append(log.TableTags, tag)
				}
			})
		})
		// TODO: Seemingly impossible to retrieve the key mode from play log?
		ForEach(out, func(rlog *dto.RivalScoreDataLogDto, _ int) {
			keys := make([]int, 7)
			for i := range keys {
				keys[i] = i
			}
			rawRandomPatternArr := random.MakeRandom(rlog.Seed, keys)
			randomPattern := ""
			for _, p := range rawRandomPatternArr {
				randomPattern += fmt.Sprintf("%d", p+1)
			}
			rlog.RandomPattern = randomPattern
		})
		return nil
	})
	return
}

// Similar to syncScoreDataLog but not delete any old content, only append new logs
func appendScoreDataLog(tx *gorm.DB, rivalScoreDatalog []*entity.RivalScoreDataLog) error {
	return tx.Model(&entity.RivalScoreDataLog{}).CreateInBatches(rivalScoreDatalog, DEFAULT_BATCH_SIZE).Error
}

// Fully delete all content from rival_score_data_log and reinsert them
func syncScoreDataLog(tx *gorm.DB, rivalScoreDataLog []*entity.RivalScoreDataLog, rivalID uint) error {
	if err := tx.Unscoped().Where("rival_id = ?", rivalID).Delete(&entity.RivalScoreDataLog{}).Error; err != nil {
		return err
	}

	return tx.CreateInBatches(&rivalScoreDataLog, DEFAULT_BATCH_SIZE).Error
}

func findRivalScoreDataLogList(tx *gorm.DB, filter *vo.RivalScoreDataLogVo) (out []*dto.RivalScoreDataLogDto, n int, err error) {
	fields := `
  rival_score_data_log.*,
  strftime("%s", rival_score_data_log.record_time) as RecordTimestamp,
  sd.title as title,
  sd.sub_title as sub_title,
  sd.artist as artist,
  sd.md5 as md5,
  sd.ID as rival_song_data_id
  `
	partial := tx.Model(&entity.RivalScoreDataLog{}).Order("rival_score_data_log.record_time desc").Select(fields)
	partial = partial.Debug().Joins("left join (select * from rival_song_data group by sha256) as sd on rival_score_data_log.sha256 = sd.sha256").Scopes(
		scopeRivalScoreDataLogFilter(filter),
	)
	if filter != nil && filter.Pagination != nil {
		partial = partial.Scopes(pagination(filter.Pagination))
		if filter.SongNameLike != nil && *filter.SongNameLike != "" {
			partial = partial.Where("sd.title like ('%' || ? || '%')", filter.SongNameLike)
		}
	}

	if err = partial.Debug().Find(&out).Error; err != nil {
		return
	}

	// pagination
	if filter != nil && filter.Pagination != nil {
		var count int64
		count, err = selectRivalScoreDataLogCount(tx, filter)
		if err != nil {
			return
		}
		filter.Pagination.PageCount = calcPageCount(count, filter.Pagination.PageSize)
	}
	return
}

func findLastRivalScoreDataLog(tx *gorm.DB, filter *vo.RivalScoreDataLogVo) (*entity.RivalScoreDataLog, error) {
	ret := entity.RivalScoreDataLog{}
	err := tx.Model(&ret).
		Scopes(scopeRivalScoreDataLogFilter(filter)).
		Order("record_time desc").
		Limit(1).
		Find(&ret).
		Error
	return &ret, err
}

// Query one user's key pressed count in a year
/*
	select sum(epg+lpg+egr+lgr+egd+lgd+ebd+lbd+epr+lpr+ems+lms) as keyp, date(record_time) as record_date
	from rival_score_data_log rsdl
	where STRFTIME('%Y', rsdl.record_time) = '2025'
	group by date(record_time)
	order by record_date desc
*/
func queryUserKeyCountInYear(tx *gorm.DB, filter *vo.RivalScoreDataLogVo) ([]*dto.KeyCountDto, int, error) {
	fields := `
    sum(epg+lpg+egr+lgr+egd+lgd+ebd+lbd+epr+lpr+ems+lms) as key_count,
    date(record_time) as record_date
  `
	var out []*dto.KeyCountDto
	if err := tx.Model(&entity.RivalScoreDataLog{}).
		Select(fields).
		Group("record_date").
		Order("record_date desc").
		Scopes(scopeRivalScoreDataLogFilter(filter)).
		Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, len(out), nil
}

func selectRivalScoreDataLogCount(tx *gorm.DB, filter *vo.RivalScoreDataLogVo) (count int64, err error) {
	partial := tx.Model(&entity.RivalScoreDataLog{}).Joins("left join (select * from rival_song_data group by sha256) as sd on rival_score_data_log.sha256 = sd.sha256").Scopes(
		scopeRivalScoreDataLogFilter(filter),
	)
	if filter != nil && filter.SongNameLike != nil && *filter.SongNameLike != "" {
		partial = partial.Where("sd.title like ('%' || ? || '%')", filter.SongNameLike)
	}
	err = partial.Debug().Count(&count).Error
	return
}

// Specialized scope for vo.RivalScoreDataLogVo
func scopeRivalScoreDataLogFilter(filter *vo.RivalScoreDataLogVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		if filter.SpecifyYear != nil {
			moved = moved.Where(`STRFTIME("%Y", rival_score_data_log.record_time) = ?`, filter.SpecifyYear)
		}
		if filter.OnlyCourseLogs {
			moved = moved.Where("length(rival_score_data_log.sha256) > 64")
		}
		if filter.NoCourseLog {
			moved = moved.Where("length(rival_score_data_log.sha256) = 64")
		}
		return moved
	}
}
