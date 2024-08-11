package loader

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/common/source"
	"github.com/Catizard/lampghost/internal/data/filter"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
	"github.com/Catizard/lampghost/internal/tui/choose"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
)

type RivalDataLoader interface {
	Interest(r *rival.RivalInfo) bool
	Load(r *rival.RivalInfo, filter null.Value[filter.Filter]) ([]*score.CommonScoreLog, error)
	// loadWithFilter(r *rival.RivalInfo, filter ???) ([]*score.CommonScoreLog, error)
}

func LoadRivalData(r *rival.RivalInfo) error {
	loader := chooseLoader(r)
	logs, err := loader.Load(r, filter.NullFilter)
	if err != nil {
		return err
	}
	r.CommonScoreLog = logs
	return nil
}

func LoadTaggedRivalData(r *rival.RivalInfo, tag *rival.RivalTag) error {
	loader := chooseLoader(r)
	// hack, which breaks on LR2 :(
	var filter filter.Filter = filter.SimpleFilter {
		WhereClause: fmt.Sprintf(" where date <= %d", tag.TimeStamp),
	}
	logs, err := loader.Load(r, null.NewValue(filter, true))
	if err != nil {
		return err
	}
	r.CommonScoreLog = logs
	return nil
}

func chooseLoader(r *rival.RivalInfo) RivalDataLoader {
	if OrajaDataLoader.Interest(r) && LR2DataLoader.Interest(r) {
		// Okay, we got a trouble
		msg := "The rival [%s] registered both LR2 file and Oraja file, you have to choose one to use"
		i := choose.OpenChooseTui([]string{source.LR2, source.Oraja}, fmt.Sprintf(msg, r.Name), false)
		if i == 0 {
			r.Prefer = null.StringFrom(source.LR2)
		} else {
			r.Prefer = null.StringFrom(source.Oraja)
		}
	} else if OrajaDataLoader.Interest(r) {
		r.Prefer = null.StringFrom(source.Oraja)
	} else if LR2DataLoader.Interest(r) {
		r.Prefer = null.StringFrom(source.LR2)
	}
	log.Infof("Rival [%s]'s prefer [%s]", r.Name, r.Prefer.String)
	if !r.Prefer.Valid {
		panic("panic: no loader")
	}
	if r.Prefer.Equal(null.StringFrom(source.LR2)) {
		return LR2DataLoader
	}
	return OrajaDataLoader
}
