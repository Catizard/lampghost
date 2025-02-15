package vo

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalSongDataVo struct {
	gorm.Model

	RivalId    uint
	Md5        string
	Sha256     string
	Title      string
	SubTitle   string
	Genre      string
	Artist     string
	SubArtist  string
	Tag        string
	BackBmp    string
	Level      int32
	Difficulty int32
	MaxBpm     int32
	MinBpm     int32
	Length     int32
	Mode       int32
	Judge      int32
	Feature    int32
	Content    int32
	Date       int64
	Favorite   int32
	AddDate    int64
	Notes      int32

	// Pagination
	Page     int
	PageSize int
}

func (rivalSongData *RivalSongDataVo) Entity() *entity.RivalSongData {
	return &entity.RivalSongData{
		Model: gorm.Model{
			ID:        rivalSongData.ID,
			CreatedAt: rivalSongData.CreatedAt,
			UpdatedAt: rivalSongData.UpdatedAt,
			DeletedAt: rivalSongData.DeletedAt,
		},
		RivalId:    rivalSongData.RivalId,
		Md5:        rivalSongData.Md5,
		Sha256:     rivalSongData.Sha256,
		Title:      rivalSongData.Title,
		SubTitle:   rivalSongData.SubTitle,
		Genre:      rivalSongData.Genre,
		Artist:     rivalSongData.Artist,
		SubArtist:  rivalSongData.SubArtist,
		Tag:        rivalSongData.Tag,
		BackBmp:    rivalSongData.BackBmp,
		Level:      rivalSongData.Level,
		Difficulty: rivalSongData.Difficulty,
		MaxBpm:     rivalSongData.MaxBpm,
		MinBpm:     rivalSongData.MinBpm,
		Length:     rivalSongData.Length,
		Mode:       rivalSongData.Mode,
		Judge:      rivalSongData.Judge,
		Feature:    rivalSongData.Feature,
		Content:    rivalSongData.Content,
		Date:       rivalSongData.Date,
		Favorite:   rivalSongData.Favorite,
		AddDate:    rivalSongData.AddDate,
		Notes:      rivalSongData.Notes,
	}
}
