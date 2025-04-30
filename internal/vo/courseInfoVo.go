package vo

import (
	"strings"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/rotisserie/eris"
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

// Some tables' courses are defined in an inner field `charts`, this function is 'pushing' them up
func (courseInfo *CourseInfoVo) pushupChartsHashField() error {
	if len(courseInfo.Charts) > 0 {
		// `charts` may provide `sha256` or `md5`
		firstChartDef := courseInfo.Charts[0]
		if firstChartDef.Md5 != "" {
			courseInfo.Md5 = make([]string, 0)
			for _, chart := range courseInfo.Charts {
				courseInfo.Md5 = append(courseInfo.Md5, chart.Md5)
			}
		} else if firstChartDef.Sha256 != "" {
			courseInfo.Sha256 = make([]string, 0)
			for _, chart := range courseInfo.Charts {
				courseInfo.Sha256 = append(courseInfo.Sha256, chart.Sha256)
			}
		} else {
			return eris.New("no sha256 or md5 provides")
		}
	}
	return nil
}
