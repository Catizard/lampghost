package dto

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type CustomDiffTableDto struct {
	gorm.Model

	Name        string
	Symbol      string
	LevelOrders string
}

func (table *CustomDiffTableDto) Entity() *entity.CustomDiffTable {
	return &entity.CustomDiffTable{
		Model: gorm.Model{
			ID:        table.ID,
			CreatedAt: table.CreatedAt,
			UpdatedAt: table.UpdatedAt,
			DeletedAt: table.DeletedAt,
		},
		Name:        table.Name,
		Symbol:      table.Symbol,
		LevelOrders: table.LevelOrders,
	}
}

func NewCustomDiffTableDto(table *entity.CustomDiffTable) *CustomDiffTableDto {
	return &CustomDiffTableDto{
		Model: gorm.Model{
			ID:        table.ID,
			CreatedAt: table.CreatedAt,
			UpdatedAt: table.UpdatedAt,
			DeletedAt: table.DeletedAt,
		},
		Name:        table.Name,
		Symbol:      table.Symbol,
		LevelOrders: table.LevelOrders,
	}
}
