package dto

import (
	"strings"

	"github.com/Catizard/lampghost_wails/internal/entity"
)

type DiffTableHeaderDto struct {
	ID                 uint
	HeaderUrl          string
	DataUrl            string
	Name               string
	OriginalUrl        *string
	Symbol             string
	LevelOrders        string
	UnjoinedLevelOrder []string
	TagColor           string
	TagTextColor       string
	NoTagBuild         *int

	Contents []*DiffTableDataDto

	// Only used in QueryLevelLayered interface
	SortedLevels         []string
	LevelLayeredContents map[string][]*DiffTableDataDto

	Level string
	// NOTE: children field should never be nil
	Children []DiffTableHeaderDto
	// clear => count
	LampCount map[int]int
	// rank => count
	RankCount map[int]int
	SongCount int
	LostCount int
}

func NewDiffTableHeaderDto(header *entity.DiffTableHeader, contents []*DiffTableDataDto) *DiffTableHeaderDto {
	return &DiffTableHeaderDto{
		ID:                 header.ID,
		HeaderUrl:          header.HeaderUrl,
		DataUrl:            header.DataUrl,
		Name:               header.Name,
		OriginalUrl:        header.OriginalUrl,
		Symbol:             header.Symbol,
		LevelOrders:        header.LevelOrders,
		UnjoinedLevelOrder: strings.Split(header.LevelOrders, ","),
		Contents:           contents,
		Children:           make([]DiffTableHeaderDto, 0),
		LampCount:          make(map[int]int),
		RankCount:          make(map[int]int),
		TagColor:           header.TagColor,
		TagTextColor:       header.TagTextColor,
		NoTagBuild:         header.NoTagBuild,
	}
}

func (header *DiffTableHeaderDto) Entity() *entity.DiffTableHeader {
	return &entity.DiffTableHeader{
		HeaderUrl:    header.HeaderUrl,
		DataUrl:      header.DataUrl,
		Name:         header.Name,
		OriginalUrl:  header.OriginalUrl,
		Symbol:       header.Symbol,
		LevelOrders:  strings.Join(header.UnjoinedLevelOrder, ","),
		TagColor:     header.TagColor,
		TagTextColor: header.TagTextColor,
		NoTagBuild:   header.NoTagBuild,
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
		ID:        id,
		Name:      name,
		Level:     level,
		LampCount: make(map[int]int),
	}
}

// For exporting usage (i.e imported by beatoraja)
type DiffTableHeaderExportDto struct {
	Name       string                       `json:"name"`
	Symbol     string                       `json:"symbol"`
	HeaderUrl  string                       `json:"header_url"`
	DataUrl    string                       `json:"data_url"`
	LevelOrder []string                     `json:"level_order"`
	Courses    [][]DiffTableCourseExportDto `json:"course"`
}

type DiffTableCourseExportDto struct {
	Name       string   `json:"name"`
	Constraint []string `json:"constraint"`
	Md5        []string `json:"md5"`
}
