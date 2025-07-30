package entity

import "gorm.io/gorm"

type SongDirectory struct {
	gorm.Model
	DirectoryPath string
	DirectoryName string
}
