package vo

import (
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalScoreLogVo struct {
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

	// Pagination
	Pagination *entity.Page
	// Extra filter field
	OnlyCourseLogs bool
	NoCourseLog    bool
	// StartRecordTime <= RecordTime <= EndRecordTime
	StartRecordTime      time.Time
	EndRecordTime        time.Time
	StartRecordTimestamp int64 // overwrite StartRecordTime
	EndRecordTimestamp   int64 // overwrite EndRecordTime
	// Clear >= MinimumClear
	MinimumClear *int32
	// year(RecordTime) == specifyYear
	SpecifyYear  *string
	SongNameLike *string
	// TODO: cannot be used in `findRivalScoreLogList`
	Sha256s []string
	// difftable_data.header_id == headerId
	// This filter field would generate a left join statement, be aware!
	HeaderID uint
}

func (rivalScoreLog *RivalScoreLogVo) Entity() *entity.RivalScoreLog {
	return &entity.RivalScoreLog{
		Model: gorm.Model{
			ID:        rivalScoreLog.ID,
			CreatedAt: rivalScoreLog.CreatedAt,
			UpdatedAt: rivalScoreLog.UpdatedAt,
			DeletedAt: rivalScoreLog.DeletedAt,
		},
		RivalId:    rivalScoreLog.RivalId,
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

func (rivalScoreLog *RivalScoreLogVo) ConvTimestamp() {
	rivalScoreLog.StartRecordTime = time.Unix(rivalScoreLog.StartRecordTimestamp, 0).Local()
	rivalScoreLog.EndRecordTime = time.Unix(rivalScoreLog.EndRecordTimestamp, 0).Local()
}
