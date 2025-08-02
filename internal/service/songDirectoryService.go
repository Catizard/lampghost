package service

import (
	"path"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/rotisserie/eris"
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
	if err := tx.Unscoped().Where("1=1").Delete(&entity.SongDirectory{}).Error; err != nil {
		return eris.Wrap(err, "delete")
	}
	directories := make([]*entity.SongDirectory, 0)
	for _, p := range songDirectories {
		directories = append(directories, &entity.SongDirectory{
			DirectoryPath: p,
			DirectoryName: path.Base(p),
		})
	}
	if err := tx.Model(&entity.SongDirectory{}).CreateInBatches(directories, DEFAULT_BATCH_SIZE).Error; err != nil {
		return eris.Wrap(err, "create")
	}
	return nil
}
