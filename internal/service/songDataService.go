package service

import (
	"github.com/Catizard/bmscanner"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type SongDataService struct {
	db *gorm.DB
}

func NewSongDataService(db *gorm.DB) *SongDataService {
	return &SongDataService{
		db: db,
	}
}

func (s *SongDataService) FindSongDataList() ([]*entity.SongData, int, error) {
	var data []*entity.SongData
	if err := s.db.Find(&data).Error; err != nil {
		return nil, 0, err
	}
	return data, len(data), nil
}

// BuildSongData build a 'songdata.db' file by scanning provided BMS
// directories, used for LR2 user
// NOTE: The build result is always located at $WORKING_DIRECTORY/songdata.db.
// This could make development easier. Also, this means that user can put a
// beatoraja generated 'songdata.db' file here to 'trick' lampghost
//
// TODO: Support updating part of the directories & database
func (s *SongDataService) BuildSongData(directories []string) error {
	models, err := bmscanner.ReadDirectories(directories...)
	if err != nil {
		return err
	}
	songData := make([]*entity.SongData, 0)
	for _, model := range models {
		songData = append(songData, entity.NewSongDataFromBMSModel(model))
	}
	return s.db.CreateInBatches(songData, DEFAULT_BATCH_SIZE).Error
}
