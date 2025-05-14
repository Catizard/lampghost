package entity

import (
	"strings"

	"github.com/Catizard/bmstable"
	"gorm.io/gorm"
)

type DiffTableHeader struct {
	gorm.Model

	HeaderUrl    string
	DataUrl      string
	Name         string
	OriginalUrl  *string
	Symbol       string
	OrderNumber  int `gorm:"default:0"`
	LevelOrders  string
	TagColor     string
	TagTextColor string
	NoTagBuild   *int `gorm:"default:0"`
}

func (DiffTableHeader) TableName() string {
	return "difftable_header"
}

// Convert external difficult table definition to internal one
func NewDiffTableHeaderFromImport(importHeader *bmstable.DifficultTable) *DiffTableHeader {
	return &DiffTableHeader{
		HeaderUrl:   importHeader.HeaderURL,
		DataUrl:     importHeader.DataURL,
		Name:        importHeader.Name,
		OriginalUrl: &importHeader.OriginalURL,
		Symbol:      importHeader.Symbol,
		LevelOrders: strings.Join(importHeader.LevelOrder, ","),
	}
}
