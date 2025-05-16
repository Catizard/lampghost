package service

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"gorm.io/gorm"
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
		if filter.RivalId != 0 {
			moved = moved.Where("rival_id = ?", filter.RivalId)
		}
		return moved
	}
}
