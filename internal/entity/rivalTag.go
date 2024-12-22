package entity

import "gorm.io/gorm"

type RivalTag struct {
	gorm.Model
	RivalId   uint
	TagName   string
	Generated bool `db:"generated"`
	Timestamp int64
}

func (RivalTag) TableName() string {
	return "rival_tag"
}
