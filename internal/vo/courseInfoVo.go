package vo

import (
	"strings"

	"github.com/Catizard/lampghost_wails/internal/entity"
)

type CourseInfoVo struct {
	Name       string   `json:"name"`
	Md5        []string `json:"md5"`
	Constraint []string `json:"constraint"`
}

func (courseInfo *CourseInfoVo) Entity() *entity.CourseInfo {
	return &entity.CourseInfo{
		Name:         courseInfo.Name,
		Md5s:         strings.Join(courseInfo.Md5, ","),
		Constranints: strings.Join(courseInfo.Constraint, ","),
	}
}
