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
	Prefer          null.String // Prefer to use LR2 or Oraja database file
	CommonScoreLog  []*score.CommonScoreLog
}

type RivalInfoService interface {
	// ---------- basic methods ----------
	FindRivalInfoList(filter RivalInfoFilter) ([]*RivalInfo, int, error)
	FindRivalInfoById(id int) (*RivalInfo, error)
	InsertRivalInfo(r *RivalInfo) error
	UpdateRivalInfo(id int, updater RivalInfoUpdater) (*RivalInfo, error) 
	DeleteRivalInfo(id int) error

	// Simple wrapper of FindRivalInfoList
	// After query, open tui app and wait user select one
	ChooseOneRival(msg string, filter RivalInfoFilter) (*RivalInfo, error)
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
	if v := f.Name; v.Valid && len(v.ValueOrZero()) > 0 {
		where = append(where, "name=:name")
	}
	if v := f.NameLike; v.Valid && len(v.ValueOrZero()) > 0 {
		where = append(where, "name like concat('%', :nameLike, '%')")
	}
	return strings.Join(where, " AND ")
}

type RivalInfoUpdater struct {
	Id null.Int `db:"id"`
	Name null.String `db:"name"`
	ScoreLogPath    null.String `db:"score_log_path"`
	SongDataPath    null.String `db:"song_data_path"`
	LR2UserDataPath null.String `db:"lr2_user_data_path"`
}

func (u *RivalInfoUpdater) MergeUpdate(r *RivalInfo) {
	if v := u.Name; v.Valid && len(v.ValueOrZero()) > 0 {
		r.Name = u.Name.ValueOrZero()
	}
	if v := u.ScoreLogPath; v.Valid && len(v.ValueOrZero()) > 0 {
		r.ScoreLogPath = v 
	}
	if v := u.SongDataPath; v.Valid && len(v.ValueOrZero()) > 0 {
		r.SongDataPath = v
	}
	if v := u.LR2UserDataPath; v.Valid && len(v.ValueOrZero()) > 0 {
		r.LR2UserDataPath = v 
	}
}

func (u *RivalInfoUpdater) GenerateSetClause() string {
	candidate := []string{}
	if v := u.Name; v.Valid && len(v.ValueOrZero()) > 0 {
		candidate = append(candidate, "name=:name")
	}
	if v := u.ScoreLogPath; v.Valid && len(v.ValueOrZero()) > 0 {
		candidate = append(candidate, "score_log_path=:score_log_path")
	}
	if v := u.SongDataPath; v.Valid && len(v.ValueOrZero()) > 0 {
		candidate = append(candidate, "song_data_path=:song_data_path")
	}
	if v := u.LR2UserDataPath; v.Valid && len(v.ValueOrZero()) > 0 {
		candidate = append(candidate, "lr2_user_data_path=:lr2_user_data_path")
	}
	if (len(candidate) > 0) {
		return strings.Join(candidate, ",")
	}
	return ""
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
