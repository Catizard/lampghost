package service

import (
	"github.com/Catizard/bmscanner"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
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
	// Try converting encoding if necessary
	if len(models) == 0 {
		return nil
	}
	any := models[0]
	detector := chardet.NewTextDetector()
	r, err := detector.DetectBest([]byte(any.Title))
	cname := ""
	if err != nil {
		log.Infof("failed to detect charset: %s, skipping auto convert steps", err)
	} else {
		cname = r.Language
	}
	log.Debugf("charset name: %s", cname)
	songData := make([]*entity.SongData, 0)
	d := japanese.ShiftJIS.NewDecoder()
	for _, model := range models {
		model.Genre = convShiftJISEncoding(d, model.Genre)
		model.Title = convShiftJISEncoding(d, model.Title)
		model.SubTitle = convShiftJISEncoding(d, model.SubTitle)
		model.Artist = convShiftJISEncoding(d, model.Artist)
		model.SubArtist = convShiftJISEncoding(d, model.SubArtist)
		songData = append(songData, entity.NewSongDataFromBMSModel(model))
	}
	log.Debugf("songdata length=%d", len(songData))
	return s.db.CreateInBatches(songData, DEFAULT_BATCH_SIZE).Error
}

// Convert a string from encoding 'ch' to utf-8
// Keep slient when having error
func convShiftJISEncoding(d *encoding.Decoder, str string) string {
	r, _, err := transform.String(d, str)
	if err != nil {
		return str
	}
	return r
}
