package rival

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/data/score"
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
	SongData        []*score.SongData
	Prefer          null.String // Prefer to use LR2 or Oraja database file
	CommonScoreLog  []*score.CommonScoreLog
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

	// Load rival's data, depends on rival's config
	LoadRivalData(r *RivalInfo) error
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
