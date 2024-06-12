package applier

import (
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/score"
)

type dheader difftable.DiffTableHeader
type ddata difftable.DiffTableData

type ScoreLogApplier interface {
	apply(logs []*score.CommonScoreLog, dth *dheader, dtd *ddata) (dheader, ddata)
}
