package service

import (
	"io"
	"strings"

	"github.com/Catizard/bmscanner"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
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
	cset := ""
	if err != nil {
		log.Infof("failed to detect charset: %s, skipping auto convert steps", err)
	} else {
		cset = r.Charset
	}
	log.Debugf("cset=%s", cset)
	songData := make([]*entity.SongData, 0)
	for _, model := range models {
		if cset != "" {
			model.Genre = convEncoding(cset, model.Genre)
			model.Title = convEncoding(cset, model.Title)
			model.SubTitle = convEncoding(cset, model.SubTitle)
			model.Artist = convEncoding(cset, model.Artist)
			model.SubArtist = convEncoding(cset, model.SubArtist)
		}
		songData = append(songData, entity.NewSongDataFromBMSModel(model))
	}
	log.Debugf("songdata length=%d", len(songData))
	return s.db.CreateInBatches(songData, DEFAULT_BATCH_SIZE).Error
}

// Convert a string from encoding 'ch' to utf-8
// Keep slient when having error
func convEncoding(ch, str string) string {
	r, err := charset.NewReader(strings.NewReader(str), ch)
	if err != nil {
		return str
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return str
	}
	return string(b)
}
