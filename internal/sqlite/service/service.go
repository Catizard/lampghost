package service

import (
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/sqlite"
)

// Exposes service directly
var DiffTableHeaderService difftable.DiffTableHeaderService
var CourseInfoService difftable.CourseInfoService
var RivalInfoService rival.RivalInfoService

func InitService(db *sqlite.DB) {
	DiffTableHeaderService = sqlite.NewDiffTableHeaderService(db)
	CourseInfoService = sqlite.NewCourseInfoService(db)
	RivalInfoService = sqlite.NewRivalInfoService(db)
}
