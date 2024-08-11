package data

import (
	"sort"
	"strconv"

	"github.com/Catizard/lampghost/internal/data/difftable"
)

// Returns sorted level list from one difftable
// Example: [sl1, sl2, sl3...]
func BuildSortedLevelList(dth *difftable.DiffTableHeader) []string {
	levels := make(map[string]interface{})
	for _, v := range dth.Data {
		levels[v.Level] = new(interface{})
	}
	if len(levels) == 0 {
		panic("tableHeader.json file corrupted, no level found")
	}

	sortedLevels := make([]string, 0)
	for level := range levels {
		sortedLevels = append(sortedLevels, level)
	}

	sort.Slice(sortedLevels, func(i, j int) bool {
		ll := sortedLevels[i]
		rr := sortedLevels[j]
		ill, errL := strconv.Atoi(ll)
		irr, errR := strconv.Atoi(rr)
		if errL == nil && errR == nil {
			return ill < irr
		}
		return ll < rr
	})
	return sortedLevels
}