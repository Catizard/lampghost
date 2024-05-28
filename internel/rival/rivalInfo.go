package rival

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/Catizard/lampghost/internel/common"
	"github.com/Catizard/lampghost/internel/score"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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
	_, err := db.NamedExec(`INSERT INTO rival_info (name, score_log_path, song_data_path) VALUES (:name,:score_log_path,:song_data_path)`, info)
	if err != nil {
		return err
	}
	return nil
}

// Query rival's info by name.
// Returns error when no result.
func QueryRivalInfo(name string) ([]RivalInfo, error) {
	db := common.OpenDB()
	defer db.Close()
	var ret []RivalInfo
	err := db.Select(&ret, "SELECT * FROM rival_info where name=?", name)
	if err == nil && len(ret) == 0 {
		err = fmt.Errorf("query doesn't return any result")
	}
	return ret, err
}
