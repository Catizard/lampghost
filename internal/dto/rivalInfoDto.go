package dto

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalInfoDto struct {
	gorm.Model

	Name             string
	Type             string
	ScoreLogPath     *string
	SongDataPath     *string
	ScoreDataLogPath *string
	ScoreDataPath    *string
	PlayCount        int
	MainUser         bool
	LockTagID        uint
	ReverseImport    int

	DiffTableHeader *DiffTableHeaderDto `gorm:"-"`

	TagName string
}

func NewRivalInfoDto(rival *entity.RivalInfo) *RivalInfoDto {
	return &RivalInfoDto{
		Model: gorm.Model{
			ID:        rival.ID,
			CreatedAt: rival.CreatedAt,
			UpdatedAt: rival.UpdatedAt,
			DeletedAt: rival.DeletedAt,
		},
		Name:             rival.Name,
		Type:             rival.Type,
		ScoreLogPath:     rival.ScoreLogPath,
		SongDataPath:     rival.SongDataPath,
		ScoreDataLogPath: rival.ScoreDataLogPath,
		ScoreDataPath:    rival.ScoreDataPath,
		PlayCount:        rival.PlayCount,
		MainUser:         rival.MainUser,
		LockTagID:        rival.LockTagID,
		ReverseImport:    rival.ReverseImport,
	}
}

func (rival *RivalInfoDto) Entity() *entity.RivalInfo {
	return &entity.RivalInfo{
		Name:             rival.Name,
		Type:             rival.Type,
		PlayCount:        rival.PlayCount,
		ScoreLogPath:     rival.ScoreLogPath,
		SongDataPath:     rival.SongDataPath,
		ScoreDataLogPath: rival.ScoreDataLogPath,
		ScoreDataPath:    rival.ScoreDataPath,
		MainUser:         rival.MainUser,
		LockTagID:        rival.LockTagID,
		ReverseImport:    rival.ReverseImport,
	}
}

func NewRivalInfoDtoWithDiffTable(rival *entity.RivalInfo, header *DiffTableHeaderDto) *RivalInfoDto {
	ret := NewRivalInfoDto(rival)
	ret.DiffTableHeader = header
	return ret
}

type BeatorajaDirectoryMeta struct {
	BeatorajaDirectoryPath string
	PlayerDirectories      []string
}
