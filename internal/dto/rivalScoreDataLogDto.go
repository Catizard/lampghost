package dto

import (
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalScoreDataLogDto struct {
	gorm.Model
	RivalId    uint
	Sha256     string `gorm:"index"`
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

	Title     string
	SubTitle  string
	Artist    string
	Md5       string
	TableTags []*DiffTableTagDto `gorm:"-"`
	// Pagination
	Page      int
	PageSize  int
	PageCount int
}

func (scoreDataLog *RivalScoreDataLogDto) NewRivalScoreDataLogDto() *RivalScoreDataLogDto {
	return &RivalScoreDataLogDto{
		RivalId:    scoreDataLog.RivalId,
		Sha256:     scoreDataLog.Sha256,
		Mode:       scoreDataLog.Mode,
		Clear:      scoreDataLog.Clear,
		RecordTime: scoreDataLog.RecordTime,
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

func (scoreDataLog *RivalScoreDataLogDto) Entity() *entity.RivalScoreDataLog {
	return &entity.RivalScoreDataLog{
		RivalId:    scoreDataLog.RivalId,
		Sha256:     scoreDataLog.Sha256,
		Mode:       scoreDataLog.Mode,
		Clear:      scoreDataLog.Clear,
		RecordTime: scoreDataLog.RecordTime,
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

type KeyCountDto struct {
	RecordDate string
	KeyCount   int
}
