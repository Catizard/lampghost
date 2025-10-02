package service

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/rotisserie/eris"
	. "github.com/samber/lo"
	"gorm.io/gorm"
)

type RivalTagService struct {
	db *gorm.DB
}

func NewRivalTagService(db *gorm.DB) *RivalTagService {
	return &RivalTagService{
		db: db,
	}
}

func (s *RivalTagService) FindRivalTagList(filter *vo.RivalTagVo) ([]*dto.RivalTagDto, int, error) {
	raw, n, err := findRivalTagList(s.db, filter)
	if err != nil {
		return nil, 0, err
	}
	if n == 0 {
		return make([]*dto.RivalTagDto, 0), 0, err
	}
	ret := make([]*dto.RivalTagDto, len(raw))
	for i := range raw {
		ret[i] = dto.NewRivalTagDto(raw[i])
	}
	return ret, n, nil
}

func (s *RivalTagService) AddRivalTag(rivalTag *vo.RivalTagVo) error {
	if rivalTag == nil {
		return fmt.Errorf("AddRivalTag: rivalTag cannot be nil")
	}
	if rivalTag.RivalId == 0 {
		return fmt.Errorf("AddRivalTag: rivalID cannot be 0")
	}
	// if you don't provide the tag name, it would be supplyed with tag time
	if rivalTag.TagName == "" {
		rivalTag.TagName = rivalTag.RecordTime.Format("2006-01-02 15:04:05")
	}
	rivalTag.Generated = false
	return s.db.Create(rivalTag.Entity()).Error
}

// Rebuild one rival's tags
//
// Implementaion Details:
//  1. All user customized tags would be kept
//  2. All generated tags would be kept as much as possible
func (s *RivalTagService) SyncRivalTag(rivalID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return syncRivalTag(tx, rivalID)
	})
}

func (s *RivalTagService) DeleteRivalTagByID(rivalTagID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		rivalTag, err := findRivalTagByID(s.db, rivalTagID)
		if err != nil {
			return err
		}
		if rivalTag.Generated {
			return fmt.Errorf("DeleteRivalTagByID: cannot delete generated tag")
		}
		return deleteRivalTagByID(tx, rivalTagID)
	})
}

func (s *RivalTagService) UpdateRivalTag(param *vo.RivalTagUpdateParam) error {
	if param == nil {
		return eris.Errorf("UpdateCustomCourse: param cannot be nil")
	}
	if param.ID == 0 {
		return eris.Errorf("UpdateCustomCourse: ID cannot be 0")
	}
	rivalTag, err := s.FindRivalTagByID(param.ID)
	if err != nil {
		return eris.Wrap(err, "query rival tag by id")
	}
	if param.RecordTimestamp != nil {
		if rivalTag.Generated {
			return eris.Errorf("Cannot update generated rival tag's record time")
		}
		param.RecordTime = time.Unix((*param.RecordTimestamp)/1000, 0)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		return updateRivalTag(tx, param)
	})
}

func (s *RivalTagService) FindRivalTagByID(ID uint) (*dto.RivalTagDto, error) {
	rawEntity, err := findRivalTagByID(s.db, ID)
	if err != nil {
		return nil, err
	}
	return dto.NewRivalTagDto(rawEntity), nil
}

// Sync one rival's tags
//
// This function should keep the old tags as much as possible because 'rival_info' table
// has a field 'lock_tag_id' which is a de facto froeign key to 'rival_tag'
func syncRivalTag(tx *gorm.DB, rivalID uint) error {
	tags, n, err := buildRivalTag(tx, rivalID)
	if err != nil {
		return eris.Wrap(err, "build rival tags")
	}
	if n == 0 {
		log.Warn("[RivalTagService] No tags have been built, skip syncing tags")
		// Small hack for clearing unexpected tags
		if count, err := selectRivalTagCount(tx, &vo.RivalTagVo{
			RivalId: rivalID,
		}); err != nil {
			return err
		} else if count > 0 {
			if err := deleteGeneratedTagsByRivalID(tx, rivalID); err != nil {
				return err
			}
		}
		return nil
	}

	prev, n, err := findRivalTagList(tx, &vo.RivalTagVo{
		RivalId:   rivalID,
		Generated: true,
	})
	if err != nil {
		return eris.Wrap(err, "find rival_tag")
	}
	if n == 0 {
		// Nothing to do
		return tx.Create(tags).Error
	}

	// Delete all previous generated tags that is no longer existsed in new tags
	deleteTagIDs := FilterMap(prev, func(prevTag *entity.RivalTag, _ int) (uint, bool) {
		if slices.IndexFunc(tags, func(newTag *entity.RivalTag) bool {
			return newTag.RecordTime.Equal(prevTag.RecordTime)
		}) == -1 {
			return prevTag.ID, true
		}
		return 0, false
	})

	// NOTE: This behaivour should not be triggered easily. Only two possible ways:
	// 	1. user provides the wrong scorelog.db
	// 	2. user changed the config 'IgnoreVariant' and rebuild the tags
	if len(deleteTagIDs) != 0 {
		log.Warnf("[RivalTagService] deleteTagIDs(len=%d) is not empty", len(deleteTagIDs))
		if err := tx.Unscoped().Delete(&entity.RivalTag{}, deleteTagIDs).Error; err != nil {
			return eris.Wrap(err, "delete rival_tag")
		}
	}

	// Insert all generated tags that is not existed in previous generated tags
	insertNewTags := Filter(tags, func(newTag *entity.RivalTag, _ int) bool {
		return slices.IndexFunc(prev, func(prevTag *entity.RivalTag) bool {
			return newTag.RecordTime.Equal(prevTag.RecordTime)
		}) == -1
	})

	if len(insertNewTags) == 0 {
		return nil
	}
	log.Debugf("[RivalTagService] insert %d new tags", len(insertNewTags))
	return tx.Create(insertNewTags).Error
}

// Build one rival's new tags, based on current data
// Currently, only 'First Clear' type tags are generated
func buildRivalTag(tx *gorm.DB, rivalID uint) ([]*entity.RivalTag, int, error) {
	courseInfoDtos, _, err := findCourseInfoList(tx, nil)
	if err != nil {
		return nil, 0, eris.Wrap(err, "findCourseInfoList")
	}

	interestHashSet := make(map[string]any)
	for _, course := range courseInfoDtos {
		interestHashSet[course.GetJoinedSha256("")] = new(any)
	}

	var minimumClear int32 = entity.Normal
	rawScorelogs, _, err := findRivalScoreLogList(tx, &vo.RivalScoreLogVo{
		RivalId:        rivalID,
		OnlyCourseLogs: true,
		MinimumClear:   &minimumClear,
	})
	if err != nil {
		return nil, 0, eris.Wrap(err, "findRivalScoreLogList")
	}
	interestScoreLogs := Filter(rawScorelogs, func(log *dto.RivalScoreLogDto, _ int) bool {
		_, ok := interestHashSet[log.Sha256]
		return ok
	})

	if len(interestScoreLogs) == 0 {
		return make([]*entity.RivalTag, 0), 0, nil
	}

	sort.Slice(interestScoreLogs, func(i int, j int) bool {
		return interestScoreLogs[i].RecordTime.Before(interestScoreLogs[j].RecordTime)
	})

	tags := make([]*entity.RivalTag, 0)
	for _, courseInfo := range courseInfoDtos {
		for _, scoreLog := range interestScoreLogs {
			if scoreLog.Sha256 != courseInfo.GetJoinedSha256("") {
				continue
			}
			scoreLogMode, err := strconv.Atoi(scoreLog.Mode)
			if err != nil {
				continue
			}
			if !courseInfo.IsMatchingScoreMode(scoreLogMode) {
				continue
			}
			fct := &entity.RivalTag{
				RivalId:    rivalID,
				TagName:    courseInfo.Name + " First Clear",
				Generated:  true,
				RecordTime: scoreLog.RecordTime,
			}
			tags = append(tags, fct)
			break
		}
	}
	return tags, len(tags), nil
}

func updateRivalTag(tx *gorm.DB, param *vo.RivalTagUpdateParam) error {
	return tx.Model(&param).Updates(param).Error
}

func findRivalTagList(tx *gorm.DB, filter *vo.RivalTagVo) ([]*entity.RivalTag, int, error) {
	partial := tx.Model(&entity.RivalTag{}).Scopes(scopeRivalTagFilter(filter))
	var out []*entity.RivalTag
	if err := partial.Find(&out).Error; err != nil {
		return nil, 0, err
	}

	if filter != nil && filter.Pagination != nil {
		count, err := selectRivalTagCount(tx, filter)
		log.Debugf("[RivalTagService] findRivalTagList: count: %d", count)
		if err != nil {
			return nil, 0, err
		}
		filter.Pagination.PageCount = calcPageCount(count, filter.Pagination.PageSize)
	}
	return out, len(out), nil
}

func findRivalTagByID(tx *gorm.DB, ID uint) (*entity.RivalTag, error) {
	tag := entity.RivalTag{}
	if err := tx.First(&tag, ID).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func selectRivalTagCount(tx *gorm.DB, filter *vo.RivalTagVo) (int64, error) {
	querying := tx.Model(&entity.RivalTag{})
	if filter != nil {
		querying = querying.Where(filter.Entity())
	}
	var count int64
	if err := querying.Debug().Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func deleteRivalTagByID(tx *gorm.DB, rivalTagID uint) error {
	return tx.Unscoped().Where("ID = ?", rivalTagID).Delete(&entity.RivalTag{}).Error
}

// Removes one rival's all generated tags
func deleteGeneratedTagsByRivalID(tx *gorm.DB, rivalID uint) error {
	return tx.Unscoped().Where("rival_id = ? and generated = true", rivalID).Delete(&entity.RivalTag{}).Error
}

func scopeRivalTagFilter(filter *vo.RivalTagVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity()).Scopes(pagination(filter.Pagination), scopeInIDs(filter.IDs))
		// Extra filter fields here
		if !filter.NoIgnoreEnabled {
			moved = moved.Where("Enabled = true")
		}
		return moved
	}
}
