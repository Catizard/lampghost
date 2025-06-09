package service

import (
	"slices"
	"strconv"
	"strings"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

var shouldIgnoreConstraintDefinition []string = []string{
	"no_good",
	"no_great",
	"no_speed",
}

// NOTE: NEVER USE MD5 AT DATA PROCESSING
type CourseInfoService struct {
	db *gorm.DB
}

func NewCourseInfoSerivce(db *gorm.DB) *CourseInfoService {
	return &CourseInfoService{
		db: db,
	}
}

func (s *CourseInfoService) FindCourseInfoList(filter *vo.CourseInfoVo) ([]*dto.CourseInfoDto, int, error) {
	return findCourseInfoList(s.db, filter)
}

func (s *CourseInfoService) FindCourseInfoListWithRival(rivalID uint) ([]*dto.CourseInfoDto, int, error) {
	rawCourses, n, err := findCourseInfoList(s.db, nil)
	if err != nil {
		return nil, 0, err
	}
	if n == 0 {
		return rawCourses, n, err
	}
	logs, _, err := findRivalScoreLogList(s.db, &vo.RivalScoreLogVo{
		RivalId:        rivalID,
		OnlyCourseLogs: true,
	})
	if err != nil {
		return nil, 0, err
	}
	mergeRivalScoreLogToCourses(rawCourses, logs)
	return rawCourses, n, nil
}

// Insert a list of courses into database
func addBatchCourseInfo(tx *gorm.DB, courseInfos []*entity.CourseInfo) error {
	return tx.CreateInBatches(courseInfos, DEFAULT_BATCH_SIZE).Error
}

func delCourseInfo(tx *gorm.DB, filter *vo.CourseInfoVo) error {
	return tx.Unscoped().Where(filter.Entity()).Delete(&entity.CourseInfo{}).Error
}

// This function returns CourseInfoDto directly due to a historical problem:
// scorelog.db only records the sha256 while most courses doesn't have, what it has is "md5"
//
// Therefore, only sha256 returned is not enough, we always have to get the md5 by sha256 so this is implemented
// in this basic query function
//
// TODO: This function forces using the main user's songdata.db to build the cache, should be implemented with a config
func findCourseInfoList(tx *gorm.DB, filter *vo.CourseInfoVo) ([]*dto.CourseInfoDto, int, error) {
	partial := tx.Model(&entity.CourseInfo{})
	if filter != nil {
		partial = partial.Where(filter.Entity())
	}
	config, err := config.ReadConfig()
	if err != nil {
		return nil, 0, eris.Wrap(err, "cannot read config")
	}
	var raw []*entity.CourseInfo
	if err = partial.Find(&raw).Error; err != nil {
		return nil, 0, err
	}
	if len(raw) == 0 {
		return make([]*dto.CourseInfoDto, 0), 0, nil
	}
	cache, err := queryDefaultSongHashCache(tx)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*dto.CourseInfoDto, 0)
	for i := range raw {
		shouldIgnore := false
		if config.IgnoreVariantCourse != 0 {
			match := false
			splitedConstraints := strings.Split(raw[i].Constraints, ",")
			for _, shouldIgnoreConstraint := range shouldIgnoreConstraintDefinition {
				if slices.Contains(splitedConstraints, shouldIgnoreConstraint) {
					match = true
					break
				}
			}
			if match {
				shouldIgnore = true
			}
		}
		if shouldIgnore {
			log.Debugf("Ignoring course: %s because user has set IgnoreVariantCourse true", raw[i].Name)
			continue
		}
		out = append(out, dto.NewCourseInfoDto(raw[i], cache))
	}
	return out, len(out), nil
}

func mergeRivalScoreLogToCourses(courses []*dto.CourseInfoDto, logs []*dto.RivalScoreLogDto) {
	for _, course := range courses {
		for _, scoreLog := range logs {
			if scoreLog.Sha256 != course.GetJoinedSha256("") {
				continue
			}
			// NOTE: we cannot handle the mode correctly here, see courseInfo.go for details
			scoreLogMode, err := strconv.Atoi(scoreLog.Mode)
			if err != nil {
				// do nothing...
				continue
			}
			if scoreLogMode/100 != course.GetConstraintMode()/100 {
				continue
			}
			course.Clear = max(course.Clear, scoreLog.Clear)
			if scoreLog.Clear > entity.Failed {
				if course.FirstClearTimestamp.IsZero() {
					course.FirstClearTimestamp = scoreLog.RecordTime
				} else if course.FirstClearTimestamp.After(scoreLog.RecordTime) {
					course.FirstClearTimestamp = scoreLog.RecordTime
				}
			}
		}
	}
}
