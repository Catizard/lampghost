package entity

import (
	"gorm.io/gorm"
)

type FolderContent struct {
	gorm.Model

	FolderID   uint
	FolderName string
	Sha256     string
	Md5        string
	Title      string
}

func (FolderContent) TableName() string {
	return "folder_content"
}

func FromSongDataToFolderContent(folder *Folder, songData *RivalSongData) *FolderContent {
	return &FolderContent{
		FolderID:   folder.ID,
		FolderName: folder.FolderName,
		Sha256:     songData.Sha256,
		Md5:        songData.Md5,
		Title:      songData.Title,
	}
}
