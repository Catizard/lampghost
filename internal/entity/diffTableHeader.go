package entity

import "gorm.io/gorm"

type DiffTableHeader struct {
	gorm.Model

	HeaderUrl   string
	DataUrl     string
	Name        string
	OriginalUrl *string
	Symbol      string
	LevelOrders string
}

func (DiffTableHeader) TableName() string {
	return "difftable_header"
}
