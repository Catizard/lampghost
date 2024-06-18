package builder

import (
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
)

type TagBuildParam struct {
	RivalInfo       *rival.RivalInfo
	DiffTableHeader []*difftable.DiffTableHeader
	SongScoreLog    []*score.CommonScoreLog
	CourseScoreLog  []*score.CommonScoreLog
	Courses         []*difftable.CourseInfo
}

// TODO: How can we generate LR2 users' tags?
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

func Build(ctxp TagBuildParam) []*rival.RivalTag {
	tags := make([]*rival.RivalTag, 0)
	for _, builder := range AvaiableTagBuilders {
		if builder.Interest(ctxp) {
			tags = append(tags, builder.Build(ctxp)...)
		}
	}
	return tags
}
