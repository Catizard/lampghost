package rival

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/Catizard/lampghost/internel/common"
	"github.com/Catizard/lampghost/internel/score"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	rivalConfigFileName = "rivalConfig.json"
	filePerm            = 0666
)

type RivalInfo struct {
	Id           int    `db:"id"`
	Name         string `db:"name"`
	ScoreLogPath string `db:"score_log_path"`
	SongDataPath string `db:"song_data_path"`
	Tags         []RivalTag
	ScoreLog     []score.ScoreLog
	SongData     []score.SongData
}

func InitRivalInfoTable() error {
	db, err := sqlx.Open("sqlite3", common.DBFileName)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("DROP TABLE IF EXISTS 'rival_info';CREATE TABLE rival_info (name TEXT(255) NOT NULL, id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, score_log_path TEXT(255) NOT NULL, song_data_path TEXT(255) NOT NULL);")
	return err
}

func (r *RivalInfo) LoadRivalScoreLog() error {
	scoreLog, err := score.ReadScoreLogFromSqlite(r.ScoreLogPath)
	if err != nil {
		return err
	}
	r.ScoreLog = scoreLog
	return nil
}

func (r *RivalInfo) LoadRivalSongData() error {
	songData, err := score.ReadSongDataFromSqlite(r.SongDataPath)
	if err != nil {
		return err
	}
	r.SongData = songData
	return nil
}

// Save one rival info(or say, meta data) to disk.
func (info *RivalInfo) SaveRivalInfo() error {
	// Sane check before saving rival
	if path.IsAbs(info.ScoreLogPath) {
		return fmt.Errorf("sorry, absolute path is not supported")
	}
	if _, err := os.Stat(info.ScoreLogPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cannot stat %s on your file system", info.ScoreLogPath)
	}
	if _, err := os.Stat(info.SongDataPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cannot stat %s on your file system", info.SongDataPath)
	}

	db := common.OpenDB()
	defer db.Close()
	_, err := db.NamedExec(`INSERT INTO rival_info (name, score_log_path, song_data_path) values (:name, :score_log_path, :song_data_path)`, info)
	if err != nil {
		return err
	}
	return nil
}

// Query rival's info by name. Only zero or one result could be match
// Promise that if error is not nil, one rival must be matched
// Warning: If error is not nil, the first result's value has no meaning
func QueryRivalInfo(name string) (RivalInfo, error) {
	// Read disk data into mermory
	arr, err := ReadRivalInfoFromDisk(rivalConfigFileName)
	if err != nil {
		return RivalInfo{}, err
	}

	for _, v := range arr {
		if v.Name == name {
			return v, nil
		}
	}
	return RivalInfo{}, fmt.Errorf("no such a rival named %s", name)
}

// Read rivals data from disk
func ReadRivalInfoFromDisk(path string) ([]RivalInfo, error) {
	// Create file if it doesn't exist
	f, err := os.OpenFile(rivalConfigFileName, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read previous data into mermory
	var prevArray []RivalInfo
	body, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &prevArray)
	if err != nil {
		return nil, err
	}
	return prevArray, nil
}
