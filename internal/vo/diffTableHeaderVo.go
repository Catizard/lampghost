package vo

import (
	"strings"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/go-viper/mapstructure/v2"
	"gorm.io/gorm"
)

type DiffTableHeaderVo struct {
	gorm.Model
	DataUrl            string  `json:"data_url"`
	Name               string  `json:"name"`
	OriginalUrl        *string `json:"original_url"`
	Symbol             string  `json:"symbol"`
	Courses            []CourseInfoVo
	HeaderUrl          string
	EnableFallbackSort int
	LevelOrders        string
	UnjoinedLevelOrder []string `json:"level_order"`

	Level           string
	RivalID         uint
	GhostRivalID    uint
	GhostRivalTagID uint
	Pagination      *entity.Page

	// Hack, see add method of difficult table for details
	RawCourses []interface{} `json:"course"`
}

func (header *DiffTableHeaderVo) Entity() *entity.DiffTableHeader {
	return &entity.DiffTableHeader{
		Model: gorm.Model{
			ID:        header.ID,
			CreatedAt: header.CreatedAt,
			UpdatedAt: header.UpdatedAt,
			DeletedAt: header.DeletedAt,
		},
		DataUrl:            header.DataUrl,
		Name:               header.Name,
		OriginalUrl:        header.OriginalUrl,
		Symbol:             header.Symbol,
		HeaderUrl:          header.HeaderUrl,
		EnableFallbackSort: header.EnableFallbackSort,
		LevelOrders:        strings.Join(header.UnjoinedLevelOrder, ","),
	}
}

// Parse field `RawCourses` into `Courses`
func (header *DiffTableHeaderVo) ParseRawCourses() error {
	if len(header.RawCourses) == 0 {
		return nil // Okay dokey
	}
	for i := range header.RawCourses {
		if innerArray, isNested := header.RawCourses[i].([]interface{}); isNested {
			for _, data := range innerArray {
				courseInfo := CourseInfoVo{}
				if err := mapstructure.Decode(data, &courseInfo); err != nil {
					return err
				}
				if err := courseInfo.pushupChartsHashField(); err != nil {
					return err
				}
				header.Courses = append(header.Courses, courseInfo)
			}
		} else {
			for _, data := range header.RawCourses {
				courseInfo := CourseInfoVo{}
				if err := mapstructure.Decode(data, &courseInfo); err != nil {
					return err
				}
				if err := courseInfo.pushupChartsHashField(); err != nil {
					return err
				}
				header.Courses = append(header.Courses, courseInfo)
			}
		}
	}

	return nil
}
