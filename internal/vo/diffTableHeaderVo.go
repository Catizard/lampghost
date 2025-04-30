package vo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"github.com/go-viper/mapstructure/v2"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"

	. "github.com/samber/lo"
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
	// TODO: DiffTableHeaderVo has two sort field, it works as follow:
	//  If SortBy is nil or empty string, sorting by SortNumber field
	//  Otherwise, sorting by combination of SortBy and SortOrder
	// As for now, this strategy isn't implemented
	SortBy *string
	// NOTE: NEVER access this field directly, use GetOrder() instead
	SortOrder *string
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

func (header *DiffTableHeaderVo) GetOrder() string {
	if header.SortOrder == nil {
		return "asc"
	}
	switch *header.SortOrder {
	case "ascend":
		return "asc"
	case "descend":
		return "desc"
	default:
		log.Warnf("unexpected SortOrder value: %s", *header.SortOrder)
		return "asc"
	}
}

// Used only in add difficult table process
// Table "発狂PMSデータベース(lv46～)" contains a level-order array
// which mixes number and string breaks the serialization
type ImportDiffTableHeaderVo struct {
	DataUrl     string  `json:"data_url"`
	Name        string  `json:"name"`
	OriginalUrl *string `json:"original_url"`
	Symbol      string  `json:"symbol"`
	HeaderUrl   string
	LevelOrders []any `json:"level_order"`
	RawCourses  []any `json:"course"`
	Courses     []CourseInfoVo
}

func (header *ImportDiffTableHeaderVo) PortBack() *DiffTableHeaderVo {
	castedLevelOrders := Map(header.LevelOrders, func(l any, _ int) string {
		if l, ok := l.(string); ok {
			return l
		}
		if l, ok := l.(int); ok {
			return strconv.Itoa(l)
		}
		return fmt.Sprintf("%v", l)
	})
	return &DiffTableHeaderVo{
		DataUrl:            header.DataUrl,
		Name:               header.Name,
		OriginalUrl:        header.OriginalUrl,
		Symbol:             header.Symbol,
		HeaderUrl:          header.HeaderUrl,
		UnjoinedLevelOrder: castedLevelOrders,
		Courses:            header.Courses,
	}
}

// Parse field 'RawCourses' into 'Courses'
//
// Possible layouts:
//  1. courses is an array of valid courses
//  2. courses is a two-dimensional array, every element might be an array of valid courses
//  3. courses is an array of a wrapped struct, the real courses are laid inside 'charts' field
//
// For Item3, see pushupChartsHashField for details
func (header *ImportDiffTableHeaderVo) ParseRawCourses() error {
	if len(header.RawCourses) == 0 {
		return nil // Okay dokey
	}
	for i := range header.RawCourses {
		if innerArray, isNested := header.RawCourses[i].([]any); isNested {
			for _, data := range innerArray {
				courseInfo := CourseInfoVo{}
				if err := mapstructure.Decode(data, &courseInfo); err != nil {
					return err
				}
				if err := courseInfo.pushupChartsHashField(); err != nil {
					return eris.Wrapf(err, "cannot pushup up charts hash field")
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
					return eris.Wrapf(err, "cannot pushup up charts hash field")
				}
				header.Courses = append(header.Courses, courseInfo)
			}
		}
	}

	return nil
}
