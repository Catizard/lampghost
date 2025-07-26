package service

import (
	"path"

	"github.com/Catizard/lampghost_wails/internal/database"
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

func (s *SongDirectoryService) FindSongDirectories() ([]*entity.SongDirectory, int, error) {
	return findSongDirectories(s.db)
}

// Save new bms directories definition
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
	// Regenerate the songdata.db file
	songDataDB, err := database.NewSelfGeneratedSongDataDatabase(true)
	if err != nil {
		return eris.Wrap(err, "open database")
	}
	songDataService := NewSongDataService(songDataDB)
	if err := songDataService.BuildSongData(songDirectories); err != nil {
		return eris.Wrap(err, "generating songdata.db")
	}
	return nil
}

func findSongDirectories(tx *gorm.DB) (out []*entity.SongDirectory, n int, err error) {
	err = tx.Model(&entity.SongDirectory{}).Find(&out).Error
	n = len(out)
	return
}
