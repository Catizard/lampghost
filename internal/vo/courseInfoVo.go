package vo

import (
	"strings"

	"github.com/Catizard/lampghost_wails/internal/entity"
)

type CourseInfoVo struct {
	Name       string        `json:"name"`
	Md5        []string      `json:"md5"`
	Sha256     []string      `json:"sha256"`
	Constraint []string      `json:"constraint"`
	Charts     []ChartInfoVo `json:"charts"`
	HeaderID   uint
}

// NOTE: This struct is only used in header's parse step
type ChartInfoVo struct {
	Title    string
	SubTitle string
	Artist   string
	Sha256   string `json:"sha256"`
	Md5      string `json:"md5"`
}

func (courseInfo *CourseInfoVo) Entity() *entity.CourseInfo {
	return &entity.CourseInfo{
		Name:        courseInfo.Name,
		Md5s:        strings.Join(courseInfo.Md5, ","),
		Sha256s:     strings.Join(courseInfo.Sha256, ","),
		Constraints: strings.Join(courseInfo.Constraint, ","),
		HeaderID:    courseInfo.HeaderID,
	}
}
