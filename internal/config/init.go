package config

import (
	"log"

	"github.com/Catizard/lampghost/internal/common"
	"github.com/Catizard/lampghost/internal/difftable"
	"github.com/Catizard/lampghost/internal/rival"
)

// Initialize lampghost application's database
// Don't return error, the caller cannot handle any error from InitLampGhost
func InitLampGhost() {
	db := common.OpenDB()
	tx := db.MustBegin()
	// difftable_header
	if err := difftable.InitDiffTableHeaderTable(tx); err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	// TODO: Should we clear any .json file too?
	// course_info
	if err := difftable.InitCourseInfoTable(tx); err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	// rival_info
	if err := rival.InitRivalInfoTable(tx); err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	// rival_tag
	if err := rival.InitRivalTagTable(tx); err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
