package vo

import (
	"strings"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type CustomDiffTableVo struct {
	gorm.Model

	Name        string
	Symbol      string
	LevelOrders string

	// Pagination
	Pagination *entity.Page
	// Extra Filters
	IgnoreDefaultTable bool // When flagged, exclude the default 'lampghost' table from the query result
	// Unused
	UnjoinedLevelOrder []string `json:"level_order"`
}

func (table *CustomDiffTableVo) Entity() *entity.CustomDiffTable {
	return &entity.CustomDiffTable{
		Model: gorm.Model{
			ID:        table.ID,
			CreatedAt: table.CreatedAt,
			UpdatedAt: table.UpdatedAt,
			DeletedAt: table.DeletedAt,
		},
		Name:        table.Name,
		Symbol:      table.Symbol,
		LevelOrders: strings.Join(table.UnjoinedLevelOrder, ","),
	}
}
