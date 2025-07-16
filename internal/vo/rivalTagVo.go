package vo

import (
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalTagVo struct {
	gorm.Model

	RivalId    uint
	TagName    string
	Generated  bool
	Enabled    bool
	RecordTime time.Time
	Symbol     string

	// Pagination
	Pagination *entity.Page

	// When flagged, ignore 'Enabled' field in filter statement
	// Defaults to false
	NoIgnoreEnabled bool

	RecordTimestamp *int64 // overwrites to RecordTime

	IDs []uint
}

func NewRivalTagVo(rivalTag *entity.RivalTag) *RivalTagVo {
	return &RivalTagVo{
		Model: gorm.Model{
			ID:        rivalTag.ID,
			CreatedAt: rivalTag.CreatedAt,
			UpdatedAt: rivalTag.UpdatedAt,
			DeletedAt: rivalTag.DeletedAt,
		},
		RivalId:    rivalTag.RivalId,
		TagName:    rivalTag.TagName,
		Generated:  rivalTag.Generated,
		Enabled:    rivalTag.Enabled,
		RecordTime: rivalTag.RecordTime,
		Symbol:     rivalTag.Symbol,
	}
}

func (rivalTag *RivalTagVo) Entity() *entity.RivalTag {
	return &entity.RivalTag{
		Model: gorm.Model{
			ID:        rivalTag.ID,
			CreatedAt: rivalTag.CreatedAt,
			UpdatedAt: rivalTag.UpdatedAt,
			DeletedAt: rivalTag.DeletedAt,
		},
		RivalId:    rivalTag.RivalId,
		TagName:    rivalTag.TagName,
		Generated:  rivalTag.Generated,
		Enabled:    rivalTag.Enabled,
		RecordTime: rivalTag.RecordTime,
		Symbol:     rivalTag.Symbol,
	}
}

type RivalTagUpdateParam struct {
	ID              uint
	TagName         *string
	Enabled         *bool
	RecordTime      time.Time
	RecordTimestamp *int64 `gorm:"-"` // Need to be converted by invoke side
	Symbol          *string
}

func (RivalTagUpdateParam) TableName() string {
	return "rival_tag"
}
