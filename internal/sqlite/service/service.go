package service

import (
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/sqlite"
)

// Exposes service directly
var DiffTableHeaderService difftable.DiffTableHeaderService

func InitService(db *sqlite.DB) {
	DiffTableHeaderService = sqlite.NewDiffTableHeaderService(db)
}
