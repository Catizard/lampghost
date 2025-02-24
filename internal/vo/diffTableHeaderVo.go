package vo

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type DiffTableHeaderVo struct {
	gorm.Model
	DataUrl     string           `json:"data_url"`
	Name        string           `json:"name"`
	OriginalUrl *string          `json:"original_url"`
	Symbol      string           `json:"symbol"`
	Courses     [][]CourseInfoVo `json:"course"`
	HeaderUrl   string

	Level           string
	RivalID         uint
	GhostRivalID    uint
	GhostRivalTagID uint
}

func (header *DiffTableHeaderVo) Entity() *entity.DiffTableHeader {
	return &entity.DiffTableHeader{
		Model: gorm.Model{
			ID:        header.ID,
			CreatedAt: header.CreatedAt,
			UpdatedAt: header.UpdatedAt,
			DeletedAt: header.DeletedAt,
		},
		DataUrl:     header.DataUrl,
		Name:        header.Name,
		OriginalUrl: header.OriginalUrl,
		Symbol:      header.Symbol,
		HeaderUrl:   header.HeaderUrl,
	}
}
