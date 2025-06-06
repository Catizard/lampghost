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
	IgnoreSha256          *string
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
		FolderName:    folder.FolderName,
		CustomTableID: folder.CustomTableID,
	}
}

type BindToFolderVo struct {
	FolderIDs []uint

	Title  string
	Md5    string
	Sha256 string
}
