package service

import (
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/service/impl"
	"github.com/Catizard/lampghost/internal/sqlite"
)

// Exposes service directly
var DiffTableHeaderService difftable.DiffTableHeaderService
var CourseInfoService difftable.CourseInfoService
var RivalInfoService rival.RivalInfoService
var RivalTagService rival.RivalTagService

func InitService(db *sqlite.DB) {
	DiffTableHeaderService = impl.NewDiffTableHeaderService(db)
	CourseInfoService = impl.NewCourseInfoService(db)
	RivalInfoService = impl.NewRivalInfoService(db)
	RivalTagService = impl.NewRivalTagService(db)
}
