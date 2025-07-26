package entity

import (
	"time"

	"gorm.io/gorm"
)

type RivalScoreLog struct {
	gorm.Model
	RivalId    uint
	Sha256     string `gorm:"index"`
	Md5        string // DON'T USE THIS FIELD
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

func FromRawLR2LogToRivalScoreLog(lr2Log *LR2Log) RivalScoreLog {
	return RivalScoreLog{
		Md5:        lr2Log.MD5,
		Mode:       "0", // TODO: Cannot figure out play mode here
		Clear:      int32(ConvLR2Clear(lr2Log.Clear)),
		RecordTime: time.Unix(int64(lr2Log.RowID), 0),
		OldClear:   0, // NOTE: LR2 doesn't provide old clear
		Combo:      int32(lr2Log.MaxCombo),
		Minbp:      int32(lr2Log.Minbp),
	}
}
