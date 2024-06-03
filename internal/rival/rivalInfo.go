package rival

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/Catizard/lampghost/internal/common"
	"github.com/Catizard/lampghost/internal/score"
	"github.com/Catizard/lampghost/internal/tui/choose"
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

func (r *RivalInfo) String() string {
	return fmt.Sprintf("%s (log=[%s],data=[%s])", r.Name, r.ScoreLogPath, r.SongDataPath)
}

// Initialize rival_info table
func InitRivalInfoTable() error {
	db, err := sqlx.Open("sqlite3", common.DBFileName)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("DROP TABLE IF EXISTS 'rival_info';CREATE TABLE rival_info (name TEXT(255) NOT NULL, id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, score_log_path TEXT(255) NOT NULL, song_data_path TEXT(255) NOT NULL);")
	return err
}

// Simple wrapper of LoadRivalScoreLog and LoadRivalSongData
// Only loads if field is nil
// If any error occurs, data on r cannot be insured
func (r *RivalInfo) LoadDataIfNil() error {
	if r.ScoreLog == nil {
		if err := r.loadRivalScoreLog(); err != nil {
			return err
		}
	}
	// TODO: support "shrink" mode
	if r.SongData == nil {
		if err := r.loadRivalSongData(); err != nil {
			return err
		}
	}
	return nil
}

// Like LoadDataIfNil, without nil check
func (r *RivalInfo) LoadDataForcely() error {
	if err := r.loadRivalScoreLog(); err != nil {
		return err
	}
	if err := r.loadRivalSongData(); err != nil {
		return err
	}
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
	if _, err := db.NamedExec(`INSERT INTO rival_info (name, score_log_path, song_data_path) VALUES (:name,:score_log_path,:song_data_path)`, info); err != nil {
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

// Simple wrapper of QueryRivalInfo.
// Returns error if no result matched.
// Open tui application when multiple results matched.
func QueryRivalInfoWithChoices(name string) (RivalInfo, error) {
	rivalArr, err := QueryRivalInfo(name)
	if err != nil {
		return RivalInfo{}, err
	}
	choices := make([]string, 0)
	for _, r := range rivalArr {
		choices = append(choices, r.String())
	}
	index := choose.OpenChooseTuiSkippable(choices, fmt.Sprintf("Multiple rivals named %s, please choose one:", name))
	return rivalArr[index], nil
}

func DeleteRivalInfo(id int) error {
	db := common.OpenDB()
	defer db.Close()
	_, err := db.Exec(`DELETE FROM rival_info where id=?`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RivalInfo) loadRivalScoreLog() error {
	scoreLog, err := score.ReadScoreLogFromSqlite(r.ScoreLogPath)
	if err != nil {
		return err
	}
	r.ScoreLog = scoreLog
	return nil
}

func (r *RivalInfo) loadRivalSongData() error {
	songData, err := score.ReadSongDataFromSqlite(r.SongDataPath)
	if err != nil {
		return err
	}
	r.SongData = songData
	return nil
}
