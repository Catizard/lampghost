package loader

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/data"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
	"github.com/Catizard/lampghost/internal/sqlite"
	"github.com/guregu/null/v5"
)

var LR2DataLoader RivalDataLoader = newLR2DataLoader()

type lr2DataLoader struct {
}

func newLR2DataLoader() *lr2DataLoader {
	return &lr2DataLoader{}
}

func (l *lr2DataLoader) Interest(r *rival.RivalInfo) bool {
	return r.LR2UserDataPath.Valid
}

func (l *lr2DataLoader) Load(r *rival.RivalInfo, filter null.Value[data.Filter]) ([]*score.CommonScoreLog, error) {
	if !l.Interest(r) {
		return nil, fmt.Errorf("[LR2DataLoader] cannot load")
	}

	// Directly read from scorelog table
	rawLogs, err := sqlite.DirectlyLoadTable[score.LR2ScoreLog](r.LR2UserDataPath.String, "score")
	if err != nil {
		return nil, err
	}

	// Convert raw data to common form
	logs := make([]*score.CommonScoreLog, 0)
	for _, rawLog := range rawLogs {
		logs = append(logs, score.NewCommonScoreLogFromLR2(rawLog))
	}

	return logs, nil
}
