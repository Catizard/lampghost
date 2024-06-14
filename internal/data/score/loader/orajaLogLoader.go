package loader

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
	"github.com/Catizard/lampghost/internal/sqlite"
)

// TODO: OrajaLogLoader is obviously not a good name if there was multiple log loader for oraja exist.
// Note: orajaLogLoader is designed to be stateless, so we can expose it directly
var OrajaLogLoader ScoreLogLoader = newOrajaLogLoader()

type orajaLogLoader struct {
}

func newOrajaLogLoader() *orajaLogLoader {
	return &orajaLogLoader{}
}

func (l *orajaLogLoader) Interest(r *rival.RivalInfo) bool {
	return r.SongDataPath.Valid && r.ScoreLogPath.Valid
}

func (l *orajaLogLoader) Load(r *rival.RivalInfo) ([]*score.CommonScoreLog, error) {
	if !l.Interest(r) {
		return nil, fmt.Errorf("[OrajaLogLoader] cannot load")
	}

	// database initialize
	db := sqlite.NewDB(r.SongDataPath.String)
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
	rows, err := tx.Queryx("SELECT * FROM scorelog")
	if err != nil {
		return nil, err
	}
	rawLogs := make([]score.ScoreLog, 0)
	for rows.Next() {
		var log score.ScoreLog
		err = rows.StructScan(&log)
		if err != nil {
			return nil, err
		}
		rawLogs = append(rawLogs, log)
	}

	// Convert raw data to common form
	logs := make([]*score.CommonScoreLog, 0)
	for _, rawLog := range rawLogs {
		logs = append(logs, score.NewCommonScoreLogFromOraja(rawLog))
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return logs, nil
}
