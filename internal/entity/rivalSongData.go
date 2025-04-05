package entity

import (
	"gorm.io/gorm"
)

type RivalSongData struct {
	gorm.Model
	RivalId    uint
	Md5        string
	Sha256     string `gorm:"index"`
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
		BackBmp:    songData.BackBmp,
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
	}
}

type SongHashCache struct {
	md5KeyCache    map[string]string
	sha256KeyCache map[string]string
}

func NewSongHashCache(md5KeyCache map[string]string, sha256KeyCache map[string]string) *SongHashCache {
	return &SongHashCache{
		md5KeyCache:    md5KeyCache,
		sha256KeyCache: sha256KeyCache,
	}
}

func (c *SongHashCache) GetSHA256(md5 string) (string, bool) {
	v, ok := c.md5KeyCache[md5]
	return v, ok
}

func (c *SongHashCache) GetMD5(sha256 string) (string, bool) {
	v, ok := c.sha256KeyCache[sha256]
	return v, ok
}
