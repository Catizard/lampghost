package dto

import (
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalScoreLogDto struct {
	gorm.Model
	RivalId    uint
	Sha256     string
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

	Md5             string
	RivalSongDataID uint
	Title           string
	SubTitle        string
	Artist          string
	TableTags       []*DiffTableTagDto `gorm:"-"`
	RecordTimestamp int64
	// Pagination
	Page      int
	PageSize  int
	PageCount int
}

func NewRivalScoreLogDto(rivalScoreLog *entity.RivalScoreLog) *RivalScoreLogDto {
	return &RivalScoreLogDto{
		Model: gorm.Model{
			ID:        rivalScoreLog.ID,
			CreatedAt: rivalScoreLog.CreatedAt,
			UpdatedAt: rivalScoreLog.UpdatedAt,
			DeletedAt: rivalScoreLog.DeletedAt,
		},
		Sha256:     rivalScoreLog.Sha256,
		Mode:       rivalScoreLog.Mode,
		Clear:      rivalScoreLog.Clear,
		OldClear:   rivalScoreLog.OldClear,
		Score:      rivalScoreLog.Score,
		OldScore:   rivalScoreLog.OldScore,
		Combo:      rivalScoreLog.Combo,
		OldCombo:   rivalScoreLog.OldCombo,
		Minbp:      rivalScoreLog.Minbp,
		OldMinbp:   rivalScoreLog.OldMinbp,
		RecordTime: rivalScoreLog.RecordTime,
	}
}

func (rivalScoreLog *RivalScoreLogDto) Entity() *entity.RivalScoreLog {
	return &entity.RivalScoreLog{
		Model: gorm.Model{
			ID:        rivalScoreLog.ID,
			CreatedAt: rivalScoreLog.CreatedAt,
			UpdatedAt: rivalScoreLog.UpdatedAt,
			DeletedAt: rivalScoreLog.DeletedAt,
		},
		Sha256:     rivalScoreLog.Sha256,
		Mode:       rivalScoreLog.Mode,
		Clear:      rivalScoreLog.Clear,
		OldClear:   rivalScoreLog.OldClear,
		Score:      rivalScoreLog.Score,
		OldScore:   rivalScoreLog.OldScore,
		Combo:      rivalScoreLog.Combo,
		OldCombo:   rivalScoreLog.OldCombo,
		Minbp:      rivalScoreLog.Minbp,
		OldMinbp:   rivalScoreLog.OldMinbp,
		RecordTime: rivalScoreLog.RecordTime,
	}
}
