package dto

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
)

type FolderDto struct {
	ID            uint
	FolderName    string
	CustomTableID uint
	OrderNumber   int

	// Only used in tree query interface
	Contents []*FolderContentDto
}

func NewFolderDto(folder *entity.Folder, contents []*FolderContentDto) *FolderDto {
	return &FolderDto{
		ID:          folder.ID,
		FolderName:  folder.FolderName,
		Contents:    contents,
		OrderNumber: folder.OrderNumber,
	}
}

func (folder *FolderDto) Entity() *entity.Folder {
	return &entity.Folder{
		FolderName: folder.FolderName,
	}
}
