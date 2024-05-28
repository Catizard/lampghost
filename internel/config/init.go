package config

import (
	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/Catizard/lampghost/internel/rival"
)

// Initialize lampghost application's database
// Don't return error, the caller cannot handle any error from InitLampGhost
func InitLampGhost() {
	// difftable_header
	if err := difftable.InitDifftableHeaderTable(); err != nil {
		panic(err)
	}
	// TODO: difftable's content
	// rival_info
	if err := rival.InitRivalInfoTable(); err != nil {
		panic(err)
	}
	// rival_tag
	if err := rival.InitRivalTagTable(); err != nil {
		panic(err)
	}
}