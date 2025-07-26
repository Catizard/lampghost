package entity

import "gorm.io/gorm"

type SongDirectory struct {
	gorm.Model
	DirectoryPath string
	DirectoryName string
}

func (SongDirectory) TableName() string {
	return "song_directory"
}
