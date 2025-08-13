package entity

import (
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

const (
	RIVAL_TYPE_BEATORAJA = "beatoraja"
	RIVAL_TYPE_LR2       = "LR2"
)

type RivalInfo struct {
	gorm.Model
	Name             string
	Type             string `gorm:"default:beatoraja"`
	ScoreLogPath     *string
	SongDataPath     *string
	ScoreDataLogPath *string
	ScoreDataPath    *string
	PlayCount        int
	MainUser         bool
	LockTagID        uint `gorm:"default:0"` // 0 means no version lock
	ReverseImport    int  `gorm:"default:0"` // 0 means no
}

func (RivalInfo) TableName() string {
	return "rival_info"
}

func (rivalInfo *RivalInfo) BeforeCreate(tx *gorm.DB) (err error) {
	if rivalInfo.Type == "" {
		return eris.New("rival_info: type cannot be empty")
	}
	return nil
}
