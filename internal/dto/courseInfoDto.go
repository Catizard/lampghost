package dto

import (
	"strings"

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
	Constranints       string
}

func NewCourseInfoDto(courseInfo *entity.CourseInfo, cache *entity.SongHashCache) *CourseInfoDto {
	ret := &CourseInfoDto{
		ID:           courseInfo.ID,
		HeaderID:     courseInfo.HeaderID,
		Name:         courseInfo.Name,
		Md5s:         courseInfo.Md5s,
		Constranints: courseInfo.Constranints,
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
