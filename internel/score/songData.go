package score

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SongData struct {
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

func ReadSongDataFromSqlite(filePath string) ([]SongData, error) {
	if !strings.HasSuffix(filePath, ".db") {
		return nil, fmt.Errorf("try reading from sqlite database file while path doesn't contains a .db suffix")
	}
	db, err := sqlx.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Queryx("select * from song")
	if err != nil {
		return nil, err
	}
	songDataArray := make([]SongData, 0)
	for rows.Next() {
		var songData SongData
		err = rows.StructScan(&songData)
		if err != nil {
			return nil, err
		}
		songDataArray = append(songDataArray, songData)
	}
	return songDataArray, nil
}
