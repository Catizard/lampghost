package service

import (
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

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

// This is a workaround since we don't have a concreate dao/repository layer
// Also we are under race-condition
var ignoreVariantCourse bool = false

// NOTE: NEVER USE MD5 AT DATA PROCESSING
type CourseInfoService struct {
	db           *gorm.DB
	mu           sync.Mutex
	configNotify <-chan any
}

func NewCourseInfoSerivce(db *gorm.DB) *CourseInfoService {
	ret := &CourseInfoService{
		db: db,
	}
	return ret
}

func (s *CourseInfoService) SubscribeConfigChanges(conf *config.ApplicationConfig, configNotify <-chan any) *CourseInfoService {
	ignoreVariantCourse = conf.IgnoreVariantCourse != 0
	s.configNotify = configNotify
	go s.listenConfigChanges()
	return s
}

func (s *CourseInfoService) listenConfigChanges() {
	for {
		<-s.configNotify
		log.Debugf("[CourseInfoService] received config change notification")
		go func() {
			log.Debugf("[CourseInfoService] updating config")
			if conf, err := config.ReadConfig(); err != nil {
				log.Errorf("cannot read config: %s", err)
			} else {
				ignoreVariantCourse = conf.IgnoreVariantCourse != 0
			}
		}()
	}
}

func (s *CourseInfoService) FindCourseInfoList(filter *vo.CourseInfoVo) ([]*dto.CourseInfoDto, int, error) {
	return findCourseInfoList(s.db, filter)
}

func (s *CourseInfoService) FindCourseInfoListWithRival(filter *vo.CourseInfoVo) ([]*dto.CourseInfoDto, int, error) {
	rawCourses, n, err := findCourseInfoList(s.db, filter)
	if err != nil {
		return nil, 0, err
	}
	if n == 0 {
		return rawCourses, n, err
	}
	logs, _, err := findRivalScoreLogList(s.db, &vo.RivalScoreLogVo{
		RivalId:        filter.RivalID,
		OnlyCourseLogs: true,
	})
	if err != nil {
		return nil, 0, err
	}
	mergeRivalScoreLogToCourses(rawCourses, logs)
	return rawCourses, n, nil
}

func (s *CourseInfoService) FindCourseInfoByID(courseID uint) (course *entity.CourseInfo, err error) {
	err = s.db.First(&course, courseID).Error
	return
}

func (s *CourseInfoService) QueryCourseSongListWithRival(filter *vo.CourseInfoVo) (songs []*dto.RivalSongDataDto, n int, err error) {
	courseInfo, err := s.FindCourseInfoByID(filter.ID)
	if err != nil {
		return
	}
	var endGhostRecordTime time.Time
	if filter.GhostRivalTagID > 0 {
		tag, err := findRivalTagByID(s.db, filter.GhostRivalTagID)
		if err != nil {
			return nil, 0, eris.Wrap(err, "failed to query rival tag by id")
		}
		endGhostRecordTime = tag.RecordTime
	}
	queryParam := &vo.RivalSongDataVo{
		RemoveDuplicate:    true,
		RivalId:            filter.RivalID,
		GhostRivalID:       filter.GhostRivalID,
		EndGhostRecordTime: endGhostRecordTime,
	}
	if courseInfo.Md5s != "" {
		queryParam.Md5s = strings.Split(courseInfo.Md5s, ",")
	} else if courseInfo.Sha256s != "" {
		queryParam.Sha256s = strings.Split(courseInfo.Sha256s, ",")
	} else {
		err = eris.Errorf("unexpected course data which doesn't have sha256 or md5")
		return
	}
	songs, n, err = findRivalSongDataListWithRival(s.db, queryParam)
	return
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
// TODO: It's kind of hard to implement pagination feature for this query function due to:
// 1. hash identity field might be sha256s or md5s
// 2. the strategy of using sql to filter variant courses(no_speed, no_good...) also not very obvious
func findCourseInfoList(tx *gorm.DB, filter *vo.CourseInfoVo) ([]*dto.CourseInfoDto, int, error) {
	var raw []*entity.CourseInfo
	if err := tx.Debug().Model(&entity.CourseInfo{}).Scopes(scopeCourseInfoFilter(filter)).Find(&raw).Error; err != nil {
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
		if ignoreVariantCourse {
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
				next := scoreLog.RecordTime.Unix()
				if course.FirstClearTimestamp == 0 || course.FirstClearTimestamp > next {
					course.FirstClearTimestamp = next
				}
			}
		}
	}
}

// Common query scope for vo.RivalInfoVo
func scopeCourseInfoFilter(filter *vo.CourseInfoVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Add extra filter here
		return moved
	}
}
