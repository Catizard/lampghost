package vo

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type FolderContentVo struct {
	gorm.Model

	FolderID   uint
	FolderName string
	Sha256     string
	Md5        string
	Title      string

	// Extra filter fields
	IDs       []uint
	FolderIDs []uint
}

func (content *FolderContentVo) Entity() *entity.FolderContent {
	return &entity.FolderContent{
		Model: gorm.Model{
			ID:        content.ID,
			CreatedAt: content.CreatedAt,
			UpdatedAt: content.UpdatedAt,
			DeletedAt: content.DeletedAt,
		},
		FolderID:   content.FolderID,
		FolderName: content.FolderName,
		Sha256:     content.Sha256,
		Md5:        content.Md5,
		Title:      content.Title,
	}
}
