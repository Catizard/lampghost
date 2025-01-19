package entity

import (
	"gorm.io/gorm"
)

type FolderContent struct {
	gorm.Model

	FolderID        uint
	DiffTableDataID uint
	FolderName      string
	Sha256          string
	Md5             string
	Title           string
}

func (FolderContent) TableName() string {
	return "folder_content"
}

func FromDiffTableDataToFolderContent(folder *Folder, diffTableData *DiffTableData) *FolderContent {
	return &FolderContent{
		FolderID:   folder.ID,
		FolderName: folder.FolderName,
		Sha256:     diffTableData.Sha256,
		Md5:        diffTableData.Md5,
		Title:      diffTableData.Title,
	}
}
