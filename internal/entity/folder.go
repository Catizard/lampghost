package entity

import "gorm.io/gorm"

type Folder struct {
	gorm.Model

	FolderName string
	// Default value is prepared for compatibility, old version of lampghost
	// doesn't have CustomTable definition.
	// Also, this design makes the CustomTable module must contains a default
	// table which id is always 1 and can never be removed. It's not a very
	// big deal but still tedious.
	CustomTableID uint `gorm:"default:1"`
}

func (Folder) TableName() string {
	return "folder"
}

func NewFolder(folderName string) Folder {
	return Folder{
		FolderName: folderName,
	}
}
