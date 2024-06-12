package loader

import (
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
)

var LR2LogLoader ScoreLogLoader = newLR2LogLoader()

type lr2LogLoader struct {
}

func newLR2LogLoader() *lr2LogLoader {
	return &lr2LogLoader{}
}

func (l *lr2LogLoader) interest(r *rival.RivalInfo) bool {
	panic("TODO: interest")
}

func (l *lr2LogLoader) load(r *rival.RivalInfo) ([]*score.CommonScoreLog, error) {
	panic("TODO: load")
}
