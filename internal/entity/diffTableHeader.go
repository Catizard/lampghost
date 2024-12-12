package entity

import "gorm.io/gorm"

type DiffTableHeader struct {
	gorm.Model

	DataUrl     string
	Name        string
	OriginalUrl *string
	Symbol      string
}

func (DiffTableHeader) TableName() string {
	return "difftable_header"
}
