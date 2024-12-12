package entity

import "gorm.io/gorm"

type DiffTableHeader struct {
	gorm.Model

	HeaderUrl   string
	DataUrl     string  `json:"data_url"`
	Name        string  `json:"name"`
	OriginalUrl *string `json:"original_url"`
	Symbol      string  `json:"symbol"`
}

func (DiffTableHeader) TableName() string {
	return "difftable_header"
}
