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

	Contents []*DiffTableDataDto

	// Only used in QueryLevelLayered interface
	SortedLevels         []string
	LevelLayeredContents map[string][]*DiffTableDataDto

	// Only used in tree query interface
	Level string
	// NOTE: children field should never be nil
	Children []DiffTableHeaderDto
}

func NewDiffTableHeaderDto(header *entity.DiffTableHeader, contents []*DiffTableDataDto) *DiffTableHeaderDto {
	return &DiffTableHeaderDto{
		ID:          header.ID,
		HeaderUrl:   header.HeaderUrl,
		DataUrl:     header.DataUrl,
		Name:        header.Name,
		OriginalUrl: header.OriginalUrl,
		Symbol:      header.Symbol,
		Contents:    contents,
		Children:    make([]DiffTableHeaderDto, 0),
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

func NewLevelLayeredDiffTableHeaderDto(header *entity.DiffTableHeader, sortedLevels []string, levelLayeredContents map[string][]*DiffTableDataDto) *DiffTableHeaderDto {
	ret := NewDiffTableHeaderDto(header, nil)
	ret.SortedLevels = sortedLevels
	ret.LevelLayeredContents = levelLayeredContents
	return ret
}

// Represents one level node from a difficult table
//
// NOTE: almost only passed fields are significant
func NewLevelChildNode(id uint, name string, level string) *DiffTableHeaderDto {
	return &DiffTableHeaderDto{
		ID:    id,
		Name:  name,
		Level: level,
	}
}
