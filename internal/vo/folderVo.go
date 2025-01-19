package vo

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type FolderVo struct {
	gorm.Model

	FolderName string
	BitIndex   int

	// Extra filter fields
	IDs []uint
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
		BitIndex:   folder.BitIndex,
	}
}
