package entity

import (
	"gorm.io/gorm"
)

type RivalScoreLog struct {
	gorm.Model
	RivalId   uint
	Sha256    string
	Mode      string
	Clear     int32
	OldClear  int32
	Score     int32
	OldScore  int32
	Combo     int32
	OldCombo  int32
	Minbp     int32
	OldMinbp  int32
	Timestamp int64 `gorm:"column:date"`
}

func (RivalScoreLog) TableName() string {
	return "rival_score_log"
}

func FromRawScoreLogToRivalScoreLog(scorelog *ScoreLog) RivalScoreLog {
	return RivalScoreLog{
		Sha256:    scorelog.Sha256,
		Mode:      scorelog.Mode,
		Clear:     scorelog.Clear,
		OldClear:  scorelog.Clear,
		Score:     scorelog.Score,
		OldScore:  scorelog.OldScore,
		Combo:     scorelog.Combo,
		OldCombo:  scorelog.OldCombo,
		Minbp:     scorelog.Minbp,
		OldMinbp:  scorelog.OldMinbp,
		Timestamp: scorelog.TimeStamp,
	}
}
