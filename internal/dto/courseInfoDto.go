package dto

import (
	"strings"
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
)

type CourseInfoDto struct {
	ID                 uint
	HeaderID           uint
	Name               string
	Md5                []string
	Md5s               string
	Sha256             []string
	Sha256s            string
	NoSepJoinedSha256s string
	Constraints        string

	Clear               int32
	FirstClearTimestamp time.Time
	Constraint          []string
}

func NewCourseInfoDto(courseInfo *entity.CourseInfo, cache *entity.SongHashCache) *CourseInfoDto {
	ret := &CourseInfoDto{
		ID:          courseInfo.ID,
		HeaderID:    courseInfo.HeaderID,
		Name:        courseInfo.Name,
		Md5s:        courseInfo.Md5s,
		Sha256s:     courseInfo.Sha256s,
		Constraints: courseInfo.Constraints,
		Constraint:  strings.Split(courseInfo.Constraints, ","),
	}
	ret.Md5 = strings.Split(ret.Md5s, ",")
	// Only apply repair steps when sha256s is absent
	if courseInfo.Sha256s == "" {
		ret.Sha256 = make([]string, 0)
		build := true
		for _, md5 := range ret.Md5 {
			sha256, ok := cache.GetSHA256(md5)
			if !ok {
				build = false
				break
			}
			ret.Sha256 = append(ret.Sha256, sha256)
		}
		if !build {
			ret.Sha256 = nil
		} else {
			ret.Sha256s = strings.Join(ret.Sha256, ",")
			ret.NoSepJoinedSha256s = strings.Join(ret.Sha256, "")
		}
	} else {
		ret.Sha256 = strings.Split(ret.Sha256s, ",")
		ret.NoSepJoinedSha256s = strings.Join(ret.Sha256, "")
	}
	return ret
}

func (courseInfo *CourseInfoDto) RepairHash(cache *entity.SongHashCache) {
	// Only apply repair steps when sha256s is absent
	if courseInfo.Sha256s != "" {
		return
	}
	courseInfo.Sha256 = make([]string, 0)
	build := true
	for _, md5 := range courseInfo.Md5 {
		sha256, ok := cache.GetSHA256(md5)
		if !ok {
			build = false
			break
		}
		courseInfo.Sha256 = append(courseInfo.Sha256, sha256)
	}
	if !build {
		courseInfo.Sha256 = nil
	} else {
		courseInfo.Sha256s = strings.Join(courseInfo.Sha256, ",")
		courseInfo.NoSepJoinedSha256s = strings.Join(courseInfo.Sha256, "")
	}
}

func (courseInfo *CourseInfoDto) GetJoinedSha256(sep string) string {
	return strings.Join(courseInfo.Sha256, sep)
}

func (courseInfo *CourseInfoDto) GetJoinedMd5(sep string) string {
	return strings.Join(courseInfo.Md5, sep)
}

type Constraint struct {
	Name   string
	Factor int
}

func NewConstraint(name string, factor int) Constraint {
	return Constraint{
		Name:   name,
		Factor: factor,
	}
}

var ConstraintsDefinition []Constraint = make([]Constraint, 0)

func init() {
	ConstraintsDefinition = append(ConstraintsDefinition, NewConstraint("no_speed", 1))
	ConstraintsDefinition = append(ConstraintsDefinition, NewConstraint("no_good", 1))
	ConstraintsDefinition = append(ConstraintsDefinition, NewConstraint("no_great", 2))
	ConstraintsDefinition = append(ConstraintsDefinition, NewConstraint("gauge_lr2", 1))
	ConstraintsDefinition = append(ConstraintsDefinition, NewConstraint("gauge_5k", 2))
	ConstraintsDefinition = append(ConstraintsDefinition, NewConstraint("gauge_7k", 3))
	ConstraintsDefinition = append(ConstraintsDefinition, NewConstraint("gauge_9k", 4))
	ConstraintsDefinition = append(ConstraintsDefinition, NewConstraint("gauge_24k", 5))
}

// TODO: we cannot find the LN mode value and option value here, therefore the comparison
// step should ignore the LN mode value and option value too.
func (courseInfo *CourseInfoDto) GetConstraintMode() int {
	hispeed := 0
	judge := 0
	gauge := 0
	constraints := strings.Split(courseInfo.Constraints, ",")
	// NOTE: there is no very good method to handle this and it's not worthy to do so
	// so I just hardcode this part
	for _, c := range constraints {
		if "no_speed" == c {
			hispeed = 1
		}
		if "no_good" == c {
			judge = 1
		}
		if "no_great" == c {
			judge = 2
		}
		if "gauge_lr2" == c {
			gauge = 1
		}
		if "gauge_5k" == c {
			gauge = 2
		}
		if "gauge_7k" == c {
			gauge = 3
		}
		if "gauge_9k" == c {
			gauge = 4
		}
		if "gauge_24k" == c {
			gauge = 5
		}
	}
	return /* (ln ? lnmode : 0) + option * 10*/ +hispeed*100 + judge*1000 + gauge*10000
}
