package entity

import (
	"time"

	"gorm.io/gorm"
)

type RivalScoreLog struct {
	gorm.Model
	RivalId    uint
	Sha256     string `gorm:"index"`
	Mode       string
	Clear      int32
	OldClear   int32
	Score      int32
	OldScore   int32
	Combo      int32
	OldCombo   int32
	Minbp      int32
	OldMinbp   int32
	RecordTime time.Time
}

func (RivalScoreLog) TableName() string {
	return "rival_score_log"
}

func FromRawScoreLogToRivalScoreLog(scorelog *ScoreLog) RivalScoreLog {
	return RivalScoreLog{
		Sha256:     scorelog.Sha256,
		Mode:       scorelog.Mode,
		Clear:      scorelog.Clear,
		OldClear:   scorelog.OldClear,
		Score:      scorelog.Score,
		OldScore:   scorelog.OldScore,
		Combo:      scorelog.Combo,
		OldCombo:   scorelog.OldCombo,
		Minbp:      scorelog.Minbp,
		OldMinbp:   scorelog.OldMinbp,
		RecordTime: time.Unix(scorelog.TimeStamp, 0).Local(),
	}
}
