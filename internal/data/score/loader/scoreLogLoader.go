package loader

import (
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
)

type ScoreLogLoader interface {
	interest(r *rival.RivalInfo) bool
	load(r *rival.RivalInfo) ([]*score.CommonScoreLog, error)
}
