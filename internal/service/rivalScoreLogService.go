package service

import (
	"fmt"

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
	scoreLogs, n, err := pageRivalScoreLogList(tx, filter)
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, 0, err
	}
	for i := range scoreLogs {
		scoreLogs[i].Page = filter.Page
		scoreLogs[i].PageSize = filter.PageSize
		scoreLogs[i].PageCount = int((count + int64(filter.PageSize) - 1) / int64(filter.PageSize))
	}
	return scoreLogs, n, nil
}

func findRivalScoreLogList(tx *gorm.DB, filter *vo.RivalScoreLogVo) ([]*dto.RivalScoreLogDto, int, error) {
	fields := `
		rival_score_log.*,
		datetime(rival_score_log.date, 'unixepoch') as record_time,
		sd.title as title,
		sd.md5 as md5,
		sd.ID as rival_song_data_id
	`
	partial := tx.Model(&entity.RivalScoreLog{}).Select(fields)
	if filter != nil {
		partial = partial.Where(filter.Entity())
		// Extra filters
		if filter.OnlyCourseLogs {
			partial = partial.Where("length(rival_score_log.sha256) > 64")
		}
		if filter.MaximumTimestamp > 0 {
			partial = partial.Where("rival_score_log.`date` <= ?", filter.MaximumTimestamp)
		}
	}
	var out []*dto.RivalScoreLogDto
	if err := partial.Joins("left join rival_song_data sd on rival_score_log.sha256 = sd.sha256").Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, len(out), nil
}

// Extend function to findRivalScoreLogList with page query parameter
// TODO: how to merge findRivalScoreLogList and pageRivalScoreLogList?
func pageRivalScoreLogList(tx *gorm.DB, filter *vo.RivalScoreLogVo) ([]*dto.RivalScoreLogDto, int, error) {
	if filter == nil {
		return nil, 0, fmt.Errorf("Cannot call page query without pagination parameter")
	}

	filter.Page = normalizePage(filter.Page)
	filter.PageSize = normalizePageSize(filter.PageSize)

	var count int64
	if err := tx.Model(&entity.RivalScoreLog{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var out []*dto.RivalScoreLogDto
	fields := `
		rival_score_log.*,
		datetime(rival_score_log.date, 'unixepoch') as record_time,
		sd.title as title,
		sd.md5 as md5,
		sd.ID as rival_song_data_id
	`
	if err := tx.Table("rival_score_log").Select(fields).Order("rival_score_log.date desc").Where(filter.Entity()).Scopes(
		pagination(filter.Page, filter.PageSize),
	).Joins("left join rival_song_data sd on rival_score_log.sha256 = sd.sha256").Scan(&out).Error; err != nil {
		return nil, 0, err
	}
	for i := range out {
		out[i].Page = filter.Page
		out[i].PageSize = filter.PageSize
		out[i].PageCount = int((count + int64(filter.PageSize) - 1) / int64(filter.PageSize))
	}
	return out, len(out), nil
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
