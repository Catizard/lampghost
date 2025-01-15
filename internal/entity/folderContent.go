package entity

import "gorm.io/gorm"

type FolderContent struct {
	gorm.Model

	Sha256 string
	Md5    string
	Title  string
}

func (FolderContent) TableName() string {
	return "folder_content"
}
