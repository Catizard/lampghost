package builder

import (
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
)

type TagBuildParam struct {
	RivalInfo *rival.RivalInfo
	DiffTableHeader []*difftable.DiffTableHeader
	CommonScoreLog []*score.CommonScoreLog
}

type TagBuilder interface {
	// Returns whether a TagBuilder can proceed or not
	Interest(ctxp TagBuildParam) bool
	// Build tags based on ctxp, which contains basic data (e.g difficult table)
	Build(ctxp TagBuildParam) []*rival.RivalTag
}

// Exposes TagBuilders here
var AvaiableTagBuilders []TagBuilder = make([]TagBuilder, 0)

func init() {
	AvaiableTagBuilders = append(AvaiableTagBuilders, NewOrajaCourseTagBuilder())
}