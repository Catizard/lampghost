package vo

import (
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalScoreDataLogVo struct {
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

	// Pagination
	Pagination *entity.Page

	// year(RecordTime) == SpecifyYear
	SpecifyYear    *string
	SongNameLike   *string
	OnlyCourseLogs bool
	NoCourseLog    bool
}

func (scoreDataLog *RivalScoreDataLogVo) Entity() *entity.RivalScoreDataLog {
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
