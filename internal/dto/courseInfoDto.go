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
		Constraints: courseInfo.Constraints,
		Constraint:  strings.Split(courseInfo.Constraints, ","),
	}
	ret.Md5 = strings.Split(ret.Md5s, ",")
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
	return ret
}

func (courseInfo *CourseInfoDto) GetJoinedSha256(sep string) string {
	return strings.Join(courseInfo.Sha256, sep)
}

func (courseInfo *CourseInfoDto) GetJoinedMd5(sep string) string {
	return strings.Join(courseInfo.Md5, sep)
}
