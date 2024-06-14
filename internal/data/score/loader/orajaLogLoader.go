package loader

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
	"github.com/Catizard/lampghost/internal/sqlite"
)

// TODO: OrajaDataLoader is obviously not a good name if there was multiple log loader for oraja exist.
// Note: orajaDataLoader is designed to be stateless, so we can expose it directly
var OrajaDataLoader RivalDataLoader = newOrajaDataLoader()

type orajaDataLoader struct {
}

func newOrajaDataLoader() *orajaDataLoader {
	return &orajaDataLoader{}
}

func (l *orajaDataLoader) Interest(r *rival.RivalInfo) bool {
	return r.SongDataPath.Valid && r.ScoreLogPath.Valid
}

// The default OrajaDataLoader loads 2 files: songdata.db and scorelog.db
// TODO: ignore the results, directly modify rival?
func (l *orajaDataLoader) Load(r *rival.RivalInfo) ([]*score.CommonScoreLog, error) {
	if !l.Interest(r) {
		return nil, fmt.Errorf("[OrajaDataLoader] cannot load")
	}

	// 1) Loads scorelog
	// Directly read from scorelog table
	rawLogs, err := sqlite.DirectlyLoadTable[score.ScoreLog](r.ScoreLogPath.String, "scorelog")
	if err != nil {
		return nil, err
	}

	// Convert raw data to common form
	logs := make([]*score.CommonScoreLog, 0)
	for _, rawLog := range rawLogs {
		logs = append(logs, score.NewCommonScoreLogFromOraja(rawLog))
	}

	// 2) Loads songdata.db
	rawSongs, err := sqlite.DirectlyLoadTable[score.SongData](r.SongDataPath.String, "song")
	if err != nil {
		return nil, err
	}
	// This is a workaround, since interface's defininition is (r) => ([]*commonlog, error)
	// There is no place for songdata to return while LR2's data form has only one field: log
	r.SongData = rawSongs

	return logs, nil
}
