package config

import (
	"github.com/Catizard/lampghost/internal/difftable"
	"github.com/Catizard/lampghost/internal/rival"
)

// Initialize lampghost application's database
// Don't return error, the caller cannot handle any error from InitLampGhost
func InitLampGhost() {
	// difftable_header
	if err := difftable.InitDiffTableHeaderTable(); err != nil {
		panic(err)
	}
	// TODO: Should we clear any .json file too?
	// course_info
	if err := difftable.InitCourseInfoTable(); err != nil {
		panic(err)
	}
	// rival_info
	if err := rival.InitRivalInfoTable(); err != nil {
		panic(err)
	}
	// rival_tag
	if err := rival.InitRivalTagTable(); err != nil {
		panic(err)
	}
}
