package entity

import (
	"time"

	"gorm.io/gorm"
)

type RivalTag struct {
	gorm.Model
	RivalId    uint
	TagName    string
	Generated  bool
	RecordTime time.Time
}

func (RivalTag) TableName() string {
	return "rival_tag"
}
