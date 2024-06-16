package builder

import (
	"github.com/Catizard/lampghost/internal/common/source"
	"github.com/Catizard/lampghost/internal/data/rival"
)

type orajaCourseTagBuilder struct {
}

func NewOrajaCourseTagBuilder() *orajaCourseTagBuilder {
	return &orajaCourseTagBuilder{}
}

func (builder *orajaCourseTagBuilder) Interest(ctxp TagBuildParam) bool {
	return ctxp.RivalInfo.Prefer.ValueOrZero() == source.Oraja
}

func (builder *orajaCourseTagBuilder) Build(ctxp TagBuildParam) []*rival.RivalTag {
	return nil
}