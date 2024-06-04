package score

import (
	"fmt"
	"strings"

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

func ReadScoreLogFromSqlite(filePath string) ([]ScoreLog, error) {
	if !strings.HasSuffix(filePath, ".db") {
		return nil, fmt.Errorf("try reading from sqlite database file while path doesn't contains a .db suffix")
	}
	db, err := sqlx.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Queryx("select * from scorelog")
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
