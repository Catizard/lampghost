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
	RecordTime time.Time
	Pagination *entity.Page
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
		RecordTime: rivalTag.RecordTime,
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
		RecordTime: rivalTag.RecordTime,
	}
}
