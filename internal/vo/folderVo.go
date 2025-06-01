package vo

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type FolderVo struct {
	gorm.Model

	FolderName    string
	CustomTableID uint

	// Extra filter fields
	IDs []uint
	// When not nil, filter out related data
	// NOTE: `RivalSongDataID` is actually converts to Sha256, if you set both of them,
	// then `Sha256` would be overwritted
	IgnoreSha256          *string
	IgnoreRivalSongDataID *uint
	RivalID               uint
}

func (folder *FolderVo) Entity() *entity.Folder {
	return &entity.Folder{
		Model: gorm.Model{
			ID:        folder.ID,
			CreatedAt: folder.CreatedAt,
			UpdatedAt: folder.UpdatedAt,
			DeletedAt: folder.DeletedAt,
		},
		FolderName: folder.FolderName,
	}
}
