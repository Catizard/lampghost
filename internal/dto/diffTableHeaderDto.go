package dto

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
)

type DiffTableHeaderDto struct {
	ID          uint
	HeaderUrl   string
	DataUrl     string
	Name        string
	OriginalUrl *string
	Symbol      string

	Contents []DiffTableDataDto

	// Only used in QueryLevelLayered interface
	SortedLevels         []string
	LevelLayeredContents map[string][]DiffTableDataDto
}

func NewDiffTableHeaderDto(header *entity.DiffTableHeader, contents []DiffTableDataDto) *DiffTableHeaderDto {
	return &DiffTableHeaderDto{
		HeaderUrl:   header.HeaderUrl,
		DataUrl:     header.DataUrl,
		Name:        header.Name,
		OriginalUrl: header.OriginalUrl,
		Symbol:      header.Symbol,
		Contents:    contents,
	}
}

func (header *DiffTableHeaderDto) Entity() *entity.DiffTableHeader {
	return &entity.DiffTableHeader{
		HeaderUrl:   header.HeaderUrl,
		DataUrl:     header.DataUrl,
		Name:        header.Name,
		OriginalUrl: header.OriginalUrl,
		Symbol:      header.Symbol,
	}
}

func NewLevelLayeredDiffTableHeaderDto(header *entity.DiffTableHeader, sortedLevels []string, levelLayeredContents map[string][]DiffTableDataDto) *DiffTableHeaderDto {
	ret := NewDiffTableHeaderDto(header, nil)
	ret.SortedLevels = sortedLevels
	ret.LevelLayeredContents = levelLayeredContents
	return ret
}
