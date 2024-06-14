package loader

import (
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
)

type RivalDataLoader interface {
	Interest(r *rival.RivalInfo) bool
	Load(r *rival.RivalInfo) ([]*score.CommonScoreLog, error)
	// loadWithFilter(r *rival.RivalInfo, filter ???) ([]*score.CommonScoreLog, error)
}
