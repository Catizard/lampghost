package vo

import (
	"strings"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
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
