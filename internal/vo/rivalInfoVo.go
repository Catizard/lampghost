package vo

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalInfoVo struct {
	gorm.Model
	Name         string
	ScoreLogPath *string
	SongDataPath *string
	PlayCount    int
	MainUser     bool

	Pagination *entity.Page
	Locale     *string // only passed at initialized phase
}

func (rivalInfo *RivalInfoVo) Entity() *entity.RivalInfo {
	return &entity.RivalInfo{
		Model: gorm.Model{
			ID:        rivalInfo.ID,
			CreatedAt: rivalInfo.CreatedAt,
			UpdatedAt: rivalInfo.UpdatedAt,
			DeletedAt: rivalInfo.DeletedAt,
		},
		Name:         rivalInfo.Name,
		ScoreLogPath: rivalInfo.ScoreLogPath,
		SongDataPath: rivalInfo.SongDataPath,
		PlayCount:    rivalInfo.PlayCount,
		MainUser:     rivalInfo.MainUser,
	}
}
