package dto

import (
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalScoreDataDto struct {
	gorm.Model
	RivalID    uint
	Sha256     string
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
	Trophy     string
	Option     int32
	Seed       int64
	Random     int32
	State      int32
}

func (scoreData *RivalScoreDataDto) Entity() *entity.RivalScoreData {
	return &entity.RivalScoreData{
		Sha256:     scoreData.Sha256,
		Mode:       scoreData.Mode,
		Clear:      scoreData.Clear,
		RecordTime: scoreData.RecordTime,
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
