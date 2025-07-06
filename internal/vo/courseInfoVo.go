package vo

import (
	"strings"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type CourseInfoVo struct {
	gorm.Model
	Name       string        `json:"name"`
	Md5        []string      `json:"md5"`
	Sha256     []string      `json:"sha256"`
	Constraint []string      `json:"constraint"`
	Charts     []ChartInfoVo `json:"charts"`
	HeaderID   uint

	RivalID         uint
	GhostRivalID    uint
	GhostRivalTagID uint
	// Currently only be used in backend
	IgnoreVariantCourse bool
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
		Model: gorm.Model{
			ID:        courseInfo.ID,
			CreatedAt: courseInfo.CreatedAt,
			UpdatedAt: courseInfo.UpdatedAt,
			DeletedAt: courseInfo.DeletedAt,
		},
		Name:        courseInfo.Name,
		Md5s:        strings.Join(courseInfo.Md5, ","),
		Sha256s:     strings.Join(courseInfo.Sha256, ","),
		Constraints: strings.Join(courseInfo.Constraint, ","),
		HeaderID:    courseInfo.HeaderID,
	}
}
