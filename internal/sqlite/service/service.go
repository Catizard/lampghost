package service

import (
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/sqlite"
)

// Exposes service directly
var DiffTableHeaderService difftable.DiffTableHeaderService
var CourseInfoService difftable.CourseInfoService

func InitService(db *sqlite.DB) {
	DiffTableHeaderService = sqlite.NewDiffTableHeaderService(db)
	CourseInfoService = sqlite.NewCourseInfoService(db)
}
