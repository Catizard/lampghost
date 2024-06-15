package rival

import (
	"fmt"
	"strings"

	"github.com/Catizard/lampghost/internal/data/score"
	"github.com/guregu/null/v5"
	_ "github.com/mattn/go-sqlite3"
)

type RivalInfo struct {
	Id              int         `db:"id"`
	Name            string      `db:"name"`
	ScoreLogPath    null.String `db:"score_log_path"`
	SongDataPath    null.String `db:"song_data_path"`
	LR2UserDataPath null.String `db:"lr2_user_data_path"`
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
	NameLike null.String `db:"nameLike"`
}

func (f *RivalInfoFilter) GenerateWhereClause() string {
	where := []string{"1=1"}
	if v := f.Id; v.Valid {
		where = append(where, "id=:id")
	}
	if v := f.Name; v.Valid {
		where = append(where, "name=:name")
	}
	if v := f.NameLike; v.Valid {
		where = append(where, "name like concat('%', :name, '%')")
	}
	return strings.Join(where, " AND ")
}

type RivalInfoUpdate struct {
	Name *string
}

func (r *RivalInfo) String() string {
	return fmt.Sprintf("%s (log=[%s],data=[%s],user=[%s])", r.Name, r.ScoreLogPath.ValueOrZero(), r.SongDataPath.ValueOrZero(), r.LR2UserDataPath.String)
}

// Set some fields to null other than blank
func (r *RivalInfo) BlankToNull() {
	if r.ScoreLogPath.Valid && r.ScoreLogPath.String == "" {
		r.ScoreLogPath.Valid = false
	}
	if r.SongDataPath.Valid && r.SongDataPath.String == "" {
		r.SongDataPath.Valid = false
	}
	if r.LR2UserDataPath.Valid && r.LR2UserDataPath.String == "" {
		r.LR2UserDataPath.Valid = false
	}
}
