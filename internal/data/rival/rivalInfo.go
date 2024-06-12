package rival

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/common"
	"github.com/Catizard/lampghost/internal/data/score"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
	_ "github.com/mattn/go-sqlite3"
)

type RivalInfo struct {
	Id              int         `db:"id"`
	Name            string      `db:"name"`
	ScoreLogPath    null.String `db:"score_log_path"`
	SongDataPath    null.String `db:"song_data_path"`
	LR2ScoreLogPath null.String `db:"lr2_score_log_path"`
	Tags            []RivalTag
	// TODO: I'm gonna to rename ScoreLog to OrajaScoreLog or something else
	ScoreLog []score.ScoreLog
	SongData []score.SongData
}

type RivalInfoService interface {
	// ---------- basic methods ----------
	FindRivalInfoList(filter RivalInfoFilter) ([]*RivalInfo, int, error)
	FindRivalInfoById(id int) (*RivalInfo, error)
	InsertRivalInfo(r *RivalInfo) error
	DeleteRivalInfo(id int) error

	// Simple wrapper of FindRivalInfoList
	// After query, open tui app and wait user select one
	ChooseOneRival(msg string, filter RivalInfoFilter) (*RivalInfo, error)
}

type RivalInfoFilter struct {
	// Filtering fields
	Id   null.Int    `db:"id"`
	Name null.String `db:"name"`
}

type RivalInfoUpdate struct {
	Name *string
}

func (r *RivalInfo) String() string {
	return fmt.Sprintf("%s (log=[%s],data=[%s])", r.Name, r.ScoreLogPath.ValueOrZero(), r.SongDataPath.ValueOrZero())
}

// Simple wrapper of LoadRivalScoreLog and LoadRivalSongData
// Only loads if field is nil
// If any error occurs, data on r cannot be insured
// TODO: Refactor its behaviour to "load data if file specified and exists, and merge them"
func (r *RivalInfo) LoadDataIfNil(tag *RivalTag) error {
	if r.ScoreLog == nil {
		filter := score.ScoreLogFilter{}
		if tag != nil {
			filter.EndTime = null.IntFrom(tag.TimeStamp)
		}
		if err := r.loadRivalScoreLog(filter); err != nil {
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

func (r *RivalInfo) loadRivalScoreLog(filter score.ScoreLogFilter) error {
	scoreLog, err := score.ReadScoreLogFromSqlite(r.ScoreLogPath.ValueOrZero(), filter)
	if err != nil {
		return err
	}
	r.ScoreLog = scoreLog
	log.Infof("loaded %d logs\n", len(r.ScoreLog))
	return nil
}

func (r *RivalInfo) loadRivalSongData() error {
	songData, err := score.ReadSongDataFromSqlite(r.SongDataPath.ValueOrZero())
	if err != nil {
		return err
	}
	r.SongData = songData
	return nil
}
