package loader

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
	"github.com/Catizard/lampghost/internal/sqlite"
)

var LR2LogLoader ScoreLogLoader = newLR2LogLoader()

type lr2LogLoader struct {
}

func newLR2LogLoader() *lr2LogLoader {
	return &lr2LogLoader{}
}

func (l *lr2LogLoader) Interest(r *rival.RivalInfo) bool {
	return r.LR2ScoreLogPath.Valid
}

func (l *lr2LogLoader) Load(r *rival.RivalInfo) ([]*score.CommonScoreLog, error) {
	if !l.Interest(r) {
		return nil, fmt.Errorf("[LR2LogLoader] cannot load")
	}

	// database initialize
	db := sqlite.NewDB(r.LR2ScoreLogPath.String)
	if err := db.Open(); err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Directly read from scorelog table
	rows, err := tx.Queryx("SELECT * FROM score")
	if err != nil {
		return nil, err
	}
	rawLogs := make([]score.LR2ScoreLog, 0)
	for rows.Next() {
		var log score.LR2ScoreLog
		err = rows.StructScan(&log)
		if err != nil {
			return nil, err
		}
		rawLogs = append(rawLogs, log)
	}

	// Convert raw data to common form
	logs := make([]*score.CommonScoreLog, 0)
	for _, rawLog := range rawLogs {
		logs = append(logs, score.NewCommonScoreLogFromLR2(rawLog))
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return logs, nil
}
