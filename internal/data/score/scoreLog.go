package score

import (
	"fmt"
	"strings"

	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type ScoreLog struct {
	Sha256    string
	Mode      string
	Clear     int32
	OldClear  int32
	Score     int32
	OldScore  int32
	Combo     int32
	OldCombo  int32
	Minbp     int32
	OldMinbp  int32
	TimeStamp int64 `db:"date"`
}

type ScoreLogFilter struct {
	BeginTime null.Int `db:"beginTime"`
	EndTime   null.Int `db:"endTime"`
}

func (f *ScoreLogFilter) GenerateWhereClause() string {
	where := []string{"1=1"}
	if v := f.BeginTime; v.Valid {
		where = append(where, "date>=:beginTime")
	}
	if v := f.EndTime; v.Valid {
		where = append(where, "date<=:endTime")
	}
	return strings.Join(where, " AND ")
}

func ReadScoreLogFromSqlite(filePath string, filter ScoreLogFilter) ([]ScoreLog, error) {
	if !strings.HasSuffix(filePath, ".db") {
		return nil, fmt.Errorf("try reading from sqlite database file while path doesn't contains a .db suffix")
	}
	db, err := sqlx.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.NamedQuery("select * from scorelog where "+filter.GenerateWhereClause(), filter)
	if err != nil {
		return nil, err
	}
	scoreLogArray := make([]ScoreLog, 0)
	for rows.Next() {
		var log ScoreLog
		err = rows.StructScan(&log)
		if err != nil {
			return nil, err
		}
		scoreLogArray = append(scoreLogArray, log)
	}
	return scoreLogArray, nil
}

type CommonScoreLog struct {
	Sha256    null.String
	Md5       null.String
	Mode      string
	Clear     int32
	OldClear  int32
	Score     int32
	OldScore  int32
	Combo     int32
	OldCombo  int32
	Minbp     int32
	OldMinbp  int32
	TimeStamp int64 `db:"date"`
}
