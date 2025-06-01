package entity

import "gorm.io/gorm"

// Custom difficult table
//
// CustomDiffTable -> DiffTableHeader
// Folder -> Level folder
// FolderContent -> DiffTableContent
type CustomDiffTable struct {
	gorm.Model

	Name        string
	Symbol      string
	LevelOrders string
}

func (CustomDiffTable) TableName() string {
	return "custom_diff_table"
}
