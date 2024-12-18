package entity

import (
	"gorm.io/gorm"
)

type RivalSongData struct {
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
	Path       string
	Folder     string
	StageFile  string
	Banner     string
	BackBmp    string
	Preview    string
	Parent     string
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
	ChartHash  string
}

func (RivalSongData) TableName() string {
	return "rival_song_data"
}
func FromRawSongDataToRivalSongData(songData *SongData) RivalSongData {
	return RivalSongData{
		Md5:        songData.Md5,
		Sha256:     songData.Sha256,
		Title:      songData.Title,
		SubTitle:   songData.SubTitle,
		Genre:      songData.Genre,
		Artist:     songData.Artist,
		SubArtist:  songData.SubArtist,
		Tag:        songData.Tag,
		Path:       songData.Path,
		Folder:     songData.Folder,
		StageFile:  songData.StageFile,
		Banner:     songData.Banner,
		BackBmp:    songData.BackBmp,
		Preview:    songData.Preview,
		Parent:     songData.Parent,
		Level:      songData.Level,
		Difficulty: songData.Difficulty,
		MaxBpm:     songData.MaxBpm,
		MinBpm:     songData.MinBpm,
		Length:     songData.Length,
		Mode:       songData.Mode,
		Judge:      songData.Judge,
		Feature:    songData.Feature,
		Content:    songData.Content,
		Date:       songData.Date,
		Favorite:   songData.Favorite,
		AddDate:    songData.AddDate,
		Notes:      songData.Notes,
		ChartHash:  songData.ChartHash,
	}
}
