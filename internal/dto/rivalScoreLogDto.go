package dto

import (
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalScoreLogDto struct {
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

	Title string
	// yyyy-mm-dd hh:mm:ss
	RecordTime string
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
		OldClear:   rivalScoreLog.Clear,
		Score:      rivalScoreLog.Score,
		OldScore:   rivalScoreLog.OldScore,
		Combo:      rivalScoreLog.Combo,
		OldCombo:   rivalScoreLog.OldCombo,
		Minbp:      rivalScoreLog.Minbp,
		OldMinbp:   rivalScoreLog.OldMinbp,
		Timestamp:  rivalScoreLog.Timestamp,
		RecordTime: time.Unix(rivalScoreLog.Timestamp, 0).Format("2006-01-02 15:04:05"),
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
		Sha256:    rivalScoreLog.Sha256,
		Mode:      rivalScoreLog.Mode,
		Clear:     rivalScoreLog.Clear,
		OldClear:  rivalScoreLog.Clear,
		Score:     rivalScoreLog.Score,
		OldScore:  rivalScoreLog.OldScore,
		Combo:     rivalScoreLog.Combo,
		OldCombo:  rivalScoreLog.OldCombo,
		Minbp:     rivalScoreLog.Minbp,
		OldMinbp:  rivalScoreLog.OldMinbp,
		Timestamp: rivalScoreLog.Timestamp,
	}
}
