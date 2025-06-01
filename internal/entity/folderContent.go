package entity

import (
	"gorm.io/gorm"
)

type FolderContent struct {
	gorm.Model

	FolderID   uint
	FolderName string
	Sha256     string
	Md5        string `gorm:"index"`
	Title      string
	Comment    string
}

func (FolderContent) TableName() string {
	return "folder_content"
}
