package entity

import (
	"time"
)

type RivalScoreData struct {
	RivalID    uint   `gorm:"primaryKey;autoIncrement:false"`
	Sha256     string `gorm:"primaryKey"`
	Md5        string // DON'T USE THIS FIELD
	Mode       string `gorm:"primaryKey"`
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
	Trophy     string
	Option     int32
	Seed       int64
	Random     int32
	State      int32
}

func (RivalScoreData) TableName() string {
	return "rival_score_data"
}

func (scoreData *RivalScoreData) GetEXScore() int32 {
	return (scoreData.Lpg+scoreData.Epg)*2 + scoreData.Lgr + scoreData.Egr
}

func (scoreData *RivalScoreData) GetAccuracy() *float32 {
	if scoreData.Notes == 0 {
		return nil
	}
	ret := float32(scoreData.GetEXScore()*50) / float32(scoreData.Notes)
	return &ret
}

func (scoreData *RivalScoreData) GetRank() *ScoreRank {
	acc := scoreData.GetAccuracy()
	if acc == nil {
		return nil
	}
	scoreRank := GetScoreRank(*acc)
	return &scoreRank
}

func FromRawScoreDataToRivalScoreData(scoreData *ScoreData) RivalScoreData {
	return RivalScoreData{
		Sha256:     scoreData.Sha256,
		Mode:       scoreData.Mode,
		Clear:      scoreData.Clear,
		RecordTime: time.Unix(scoreData.TimeStamp, 0).Local(),
		Epg:        scoreData.Epg,
		Lpg:        scoreData.Lpg,
		Egr:        scoreData.Egr,
		Lgr:        scoreData.Lgr,
		Egd:        scoreData.Egd,
		Lgd:        scoreData.Lgd,
		Ebd:        scoreData.Ebd,
		Lbd:        scoreData.Lbd,
		Epr:        scoreData.Epr,
		Lpr:        scoreData.Lpr,
		Ems:        scoreData.Ems,
		Lms:        scoreData.Lms,
		Notes:      scoreData.Notes,
		Combo:      scoreData.Combo,
		Minbp:      scoreData.Minbp,
		PlayCount:  scoreData.PlayCount,
		ClearCount: scoreData.ClearCount,
		Trophy:     scoreData.Trophy,
		Option:     scoreData.Option,
		Seed:       scoreData.Seed,
		Random:     scoreData.Random,
		State:      scoreData.State,
	}
}
