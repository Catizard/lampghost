package entity

import "gorm.io/gorm"

const MAX_FOLDER_COUNT = 25
const BEGIN_FOLDER_INDEX = 5

type Folder struct {
	gorm.Model

	FolderName string
	BitIndex   int
}

func (Folder) Tablename() string {
	return "folder"
}
