package entity

import "gorm.io/gorm"

type Folder struct {
	gorm.Model

	FolderName string
	BitIndex   int
}

func (Folder) TableName() string {
	return "folder"
}

func NewFolder(folderName string, bitIndex int) Folder {
	return Folder{
		FolderName: folderName,
		BitIndex: bitIndex,
	}
}
