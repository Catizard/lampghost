package entity

import "gorm.io/gorm"

type RivalTag struct {
	gorm.Model
	RivalId   uint
	TagName   string
	Generated bool
	Timestamp int64
}

func (RivalTag) TableName() string {
	return "rival_tag"
}
