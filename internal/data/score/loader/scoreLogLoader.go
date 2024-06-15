package loader

import (
	"github.com/Catizard/lampghost/internal/data"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
	"github.com/guregu/null/v5"
)

type RivalDataLoader interface {
	Interest(r *rival.RivalInfo) bool
	Load(r *rival.RivalInfo, filter null.Value[data.Filter]) ([]*score.CommonScoreLog, error)
	// loadWithFilter(r *rival.RivalInfo, filter ???) ([]*score.CommonScoreLog, error)
}
