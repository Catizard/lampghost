package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type SongDirectoryService struct {
	db *gorm.DB
}

func NewSongDirectoryService(db *gorm.DB) *SongDirectoryService {
	return &SongDirectoryService{
		db: db,
	}
}

func (s *SongDirectoryService) SaveSongDirectories(songDirectories []string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return saveSongDirectories(tx, songDirectories)
	})
}

func saveSongDirectories(tx *gorm.DB, songDirectories []string) error {
	if err := tx.Unscoped().Delete(&entity.SongDirectory{}).Error; err != nil {
		return err
	}
	if err := tx.CreateInBatches(songDirectories, DEFAULT_BATCH_SIZE).Error; err != nil {
		return err
	}
	return nil
}
