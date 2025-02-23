package dto

import (
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalTagDto struct {
	gorm.Model

	RivalId   uint
	TagName   string
	Generated bool
	Timestamp int64
	TagTime   time.Time
}

func NewRivalTagDto(rivalTag *entity.RivalTag) *RivalTagDto {
	return &RivalTagDto{
		Model: gorm.Model{
			ID:        rivalTag.ID,
			CreatedAt: rivalTag.CreatedAt,
			UpdatedAt: rivalTag.UpdatedAt,
			DeletedAt: rivalTag.DeletedAt,
		},
		RivalId:   rivalTag.RivalId,
		TagName:   rivalTag.TagName,
		Generated: rivalTag.Generated,
		Timestamp: rivalTag.Timestamp,
		TagTime:   time.Unix(rivalTag.Timestamp, 0),
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
		RivalId:   rivalTag.RivalId,
		TagName:   rivalTag.TagName,
		Generated: rivalTag.Generated,
		Timestamp: rivalTag.Timestamp,
	}
}
