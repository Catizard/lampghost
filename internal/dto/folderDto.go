package dto

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/entity"
)

type FolderDto struct {
	ID         uint
	FolderName string

	// Only used in tree query interface
	Contents []FolderContentDto
}

func NewFolderDto(folder *entity.Folder, contents []FolderContentDto) *FolderDto {
	return &FolderDto{
		ID:         folder.ID,
		FolderName: folder.FolderName,
		Contents:   contents,
	}
}

func (folder *FolderDto) Entity() *entity.Folder {
	return &entity.Folder{
		FolderName: folder.FolderName,
	}
}

// Represents one folder definition in beatoraja/folder/default.json
type FolderDefinitionDto struct {
	Name string `json:"name"`
	Sql  string `json:"sql"`
}

func NewFolderDefinitionDto(folder *entity.Folder) *FolderDefinitionDto {
	return &FolderDefinitionDto{
		Name: folder.FolderName,
		Sql:  fmt.Sprintf("favorite & %d != 0", 1<<folder.BitIndex),
	}
}
