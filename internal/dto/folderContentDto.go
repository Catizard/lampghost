package dto

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
)

type FolderContentDto struct {
	ID uint

	FolderID   uint
	FolderName string
	Sha256     string
	Md5        string
	Title      string
	Comment    string

	Lamp                int
	BestRecordTimestamp int64
	TableTags           []*DiffTableTagDto `gorm:"-"`
	// Pagination
	Page      int
	PageSize  int
	PageCount int
}

func NewFolderContentDto(folderContent *entity.FolderContent) *FolderContentDto {
	return &FolderContentDto{
		ID:         folderContent.ID,
		FolderID:   folderContent.FolderID,
		FolderName: folderContent.FolderName,
		Sha256:     folderContent.Sha256,
		Md5:        folderContent.Md5,
		Title:      folderContent.Title,
		Comment:    folderContent.Comment,
	}
}

func (folderContent *FolderContentDto) Entity() *entity.FolderContent {
	return &entity.FolderContent{
		FolderID:   folderContent.FolderID,
		FolderName: folderContent.FolderName,
		Sha256:     folderContent.Sha256,
		Md5:        folderContent.Md5,
		Title:      folderContent.Title,
		Comment:    folderContent.Comment,
	}
}

type FolderContentDefinitionDto struct {
	Sha256 string
	Mask   int
}
