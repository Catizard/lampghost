package entity

import (
	"time"

	"gorm.io/gorm"
)

type RivalScoreDataLog struct {
	gorm.Model
	RivalId    uint
	Sha256     string `gorm:"index"`
	Md5        string // DON'T USE THIS FIELD
	Mode       string
	Clear      int32
	RecordTime time.Time
	Epg        int32
	Lpg        int32
	Egr        int32
	Lgr        int32
	Egd        int32
	Lgd        int32
	Ebd        int32
	Lbd        int32
	Epr        int32
	Lpr        int32
	Ems        int32
	Lms        int32
	Notes      int32
	Combo      int32
	Minbp      int32
	PlayCount  int32
	ClearCount int32
	Option     int32
	Seed       int64
	Random     int32
	State      int32
}

func (RivalScoreDataLog) TableName() string {
	return "rival_score_data_log"
}

func FromRawScoreDataLogToRivalScoreDataLog(scoreDataLog *ScoreDataLog) RivalScoreDataLog {
	return RivalScoreDataLog{
		Sha256:     scoreDataLog.Sha256,
		Mode:       scoreDataLog.Mode,
		Clear:      scoreDataLog.Clear,
		RecordTime: time.Unix(scoreDataLog.TimeStamp, 0).Local(),
		Epg:        scoreDataLog.Epg,
		Lpg:        scoreDataLog.Lpg,
		Egr:        scoreDataLog.Egr,
		Lgr:        scoreDataLog.Lgr,
		Egd:        scoreDataLog.Egd,
		Lgd:        scoreDataLog.Lgd,
		Ebd:        scoreDataLog.Ebd,
		Lbd:        scoreDataLog.Lbd,
		Epr:        scoreDataLog.Epr,
		Lpr:        scoreDataLog.Lpr,
		Ems:        scoreDataLog.Ems,
		Lms:        scoreDataLog.Lms,
		Notes:      scoreDataLog.Notes,
		Combo:      scoreDataLog.Combo,
		Minbp:      scoreDataLog.Minbp,
		PlayCount:  scoreDataLog.PlayCount,
		ClearCount: scoreDataLog.ClearCount,
		Option:     scoreDataLog.Option,
		Seed:       scoreDataLog.Seed,
		Random:     scoreDataLog.Random,
		State:      scoreDataLog.State,
	}
}

func FromRawLR2LogToRivalScoreDataLog(lr2Log *LR2Log) RivalScoreDataLog {
	return RivalScoreDataLog{
		Md5:        lr2Log.MD5,
		Mode:       "7", // TODO: Cannot figure out play mode here
		Clear:      int32(ConvLR2Clear(lr2Log.Clear)),
		RecordTime: time.Unix(int64(lr2Log.RowID), 0),
		Epg:        int32(lr2Log.Perfect),
		Egr:        int32(lr2Log.Great),
		Egd:        int32(lr2Log.Good),
		Ebd:        int32(lr2Log.Bad),
		Epr:        int32(lr2Log.Poor),
		Notes:      int32(lr2Log.TotalNotes),
		Combo:      int32(lr2Log.MaxCombo),
		Minbp:      int32(lr2Log.Minbp),
		PlayCount:  int32(lr2Log.PlayCount),
		ClearCount: int32(lr2Log.ClearCount),
	}
}

func FromRawLR2LogToRivalScoreData(lr2Log *LR2Log) RivalScoreData {
	return RivalScoreData{
		Md5:        lr2Log.MD5,
		Mode:       "7", // TODO: Cannot figure out play mode here
		Clear:      int32(ConvLR2Clear(lr2Log.Clear)),
		RecordTime: time.Unix(int64(lr2Log.RowID), 0),
		Epg:        int32(lr2Log.Perfect),
		Egr:        int32(lr2Log.Great),
		Egd:        int32(lr2Log.Good),
		Ebd:        int32(lr2Log.Bad),
		Epr:        int32(lr2Log.Poor),
		Notes:      int32(lr2Log.TotalNotes),
		Combo:      int32(lr2Log.MaxCombo),
		Minbp:      int32(lr2Log.Minbp),
		PlayCount:  int32(lr2Log.PlayCount),
		ClearCount: int32(lr2Log.ClearCount),
	}
}
