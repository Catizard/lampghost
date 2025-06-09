package entity

import (
	"gorm.io/gorm"
)

type RivalInfo struct {
	gorm.Model
	Name             string
	ScoreLogPath     *string
	SongDataPath     *string
	ScoreDataLogPath *string
	PlayCount        int
	MainUser         bool
	LockTagID        uint `gorm:"default:0"` // 0 means no version lock
	ReverseImport    int  `gorm:"default:0"` // 0 means no
}

func (RivalInfo) TableName() string {
	return "rival_info"
}
