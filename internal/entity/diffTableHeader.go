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
	// defaults to false, if flagged and LevelOrders is empty, sort levels by custom sort strategy, otherwise do nothing
	EnableFallbackSort int `gorm:"default:0"`
}

func (DiffTableHeader) TableName() string {
	return "difftable_header"
}
