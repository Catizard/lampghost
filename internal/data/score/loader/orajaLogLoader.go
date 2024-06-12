package loader

import (
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
)

// TODO: OrajaLogLoader is obviously not a good name if there was multiple log loader for oraja exist.
// Note: orajaLogLoader is designed to be stateless, so we can expose it directly
var OrajaLogLoader ScoreLogLoader = newOrajaLogLoader()

type orajaLogLoader struct {
}

func newOrajaLogLoader() *orajaLogLoader {
	return &orajaLogLoader{}
}

func (l *orajaLogLoader) interest(r *rival.RivalInfo) bool {
	panic("TODO: interest")
}

func (l *orajaLogLoader) load(r *rival.RivalInfo) ([]*score.CommonScoreLog, error) {
	panic("TODO: load")
}
