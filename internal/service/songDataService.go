package service

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
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

func (s *SongDataService) SyncFolderContentDefinition(defintion []dto.FolderContentDefinitionDto) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1) Flush all favorites
		if err := s.db.Model(&entity.SongData{}).Where("1=1").Update("favorite", 0).Error; err != nil {
			return err
		}
		// 2) rewrite all definition
		// TODO: feel bad about for-loop update
		for _, def := range defintion {
			if err := s.db.Debug().Model(&entity.SongData{}).Where("sha256 = ?", def.Sha256).Update("favorite", def.Mask).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
