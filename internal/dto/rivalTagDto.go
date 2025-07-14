package dto

import (
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalTagDto struct {
	gorm.Model

	RivalId    uint
	TagName    string
	Generated  bool
	Enabled    bool
	RecordTime time.Time
	Symbol     string

	RecordTimestamp int64
}

func NewRivalTagDto(rivalTag *entity.RivalTag) *RivalTagDto {
	return &RivalTagDto{
		Model: gorm.Model{
			ID:        rivalTag.ID,
			CreatedAt: rivalTag.CreatedAt,
			UpdatedAt: rivalTag.UpdatedAt,
			DeletedAt: rivalTag.DeletedAt,
		},
		RivalId:         rivalTag.RivalId,
		TagName:         rivalTag.TagName,
		Generated:       rivalTag.Generated,
		Enabled:         rivalTag.Enabled,
		RecordTime:      rivalTag.RecordTime,
		Symbol:          rivalTag.Symbol,
		RecordTimestamp: rivalTag.RecordTime.UnixMilli(),
	}
}

func (rivalTag *RivalTagDto) Entity() *entity.RivalTag {
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
