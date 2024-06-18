package builder

import (
	"github.com/Catizard/lampghost/internal/common/clearType"
	"github.com/Catizard/lampghost/internal/common/source"
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
)

type orajaCourseTagBuilder struct {
}

func NewOrajaCourseTagBuilder() *orajaCourseTagBuilder {
	return &orajaCourseTagBuilder{}
}

func (builder *orajaCourseTagBuilder) Interest(ctxp TagBuildParam) bool {
	return ctxp.RivalInfo.Prefer.ValueOrZero() == source.Oraja
}

// Build two kinds of tags: First Clear and First Hard Clear
func (builder *orajaCourseTagBuilder) Build(ctxp TagBuildParam) []*rival.RivalTag {
	md5MapsToCourse := make(map[string]*difftable.CourseInfo)
	for _, v := range ctxp.Courses {
		md5MapsToCourse[v.Md5s] = v
	}
	md5MapsToScoreLog := make(map[string][]score.CommonScoreLog)
	for _, v := range ctxp.CourseScoreLog {
		md5 := v.Md5.ValueOrZero()
		// Skip
		if _, ok := md5MapsToCourse[md5]; !ok {
			continue
		}
		if _, ok := md5MapsToScoreLog[md5]; !ok {
			md5MapsToScoreLog[md5] = make([]score.CommonScoreLog, 0)
		}
		md5MapsToScoreLog[md5] = append(md5MapsToScoreLog[md5], *v)
	}
	tags := make([]*rival.RivalTag, 0)
	// First Clear Tag
	for _, course := range ctxp.Courses {
		if logs, ok := md5MapsToScoreLog[course.Md5s]; !ok {
			continue // No record, continue
		} else {
			for _, log := range logs {
				if log.Clear >= clearType.Normal {
					fct := rival.RivalTag{
						TagName:   course.Name + " First Clear",
						Generated: true,
						TimeStamp: log.TimeStamp.Int64,
						TagSource: ctxp.RivalInfo.Prefer.ValueOrZero(),
					}
					tags = append(tags, &fct)
					break
				}
			}
		}
	}
	// First Hard Clear Tag
	for _, course := range ctxp.Courses {
		if logs, ok := md5MapsToScoreLog[course.Md5s]; !ok {
			continue // No record, continue
		} else {
			for _, log := range logs {
				if log.Clear >= clearType.Hard {
					fct := rival.RivalTag{
						TagName:   course.Name + " First Hard Clear",
						Generated: true,
						TimeStamp: log.TimeStamp.Int64,
						TagSource: ctxp.RivalInfo.Prefer.ValueOrZero(),
					}
					tags = append(tags, &fct)
					break
				}
			}
		}
	}
	return tags
}
