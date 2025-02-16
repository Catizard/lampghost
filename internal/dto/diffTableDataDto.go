package dto

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
)

type DiffTableDataDto struct {
	ID        uint
	HeaderID  uint
	Artist    string
	Comment   string
	Level     string
	Lr2BmsId  string
	Md5       string
	NameDiff  string
	Title     string
	Url       string
	UrlDiff   string
	Sha256    string
	Lamp      int
	GhostLamp int

	PlayCount int
}

func NewDiffTableDataDto(data *entity.DiffTableData) *DiffTableDataDto {
	return &DiffTableDataDto{
		ID:       data.ID,
		HeaderID: data.HeaderID,
		Artist:   data.Artist,
		Comment:  data.Comment,
		Level:    data.Level,
		Lr2BmsId: data.Lr2BmsId,
		Md5:      data.Md5,
		NameDiff: data.NameDiff,
		Title:    data.Title,
		Url:      data.Url,
		UrlDiff:  data.UrlDiff,
		Sha256:   data.Sha256,
	}
}

// Extends to NewDiffTableDataDto, calling RepairHash internally
func NewDiffTableDataDtoWithCache(data *entity.DiffTableData, cache *entity.SongHashCache) *DiffTableDataDto {
	ret := NewDiffTableDataDto(data)
	ret.RepairHash(cache)
	return ret
}

// []*entity.DiffTableData -> []*dto.DiffTableData
func NewDiffTableDataDtoArray(arr []*entity.DiffTableData) []*DiffTableDataDto {
	contents := make([]*DiffTableDataDto, len(arr))
	for i, rawContent := range arr {
		contents[i] = NewDiffTableDataDto(rawContent)
	}
	return contents
}

func (data *DiffTableDataDto) RepairHash(cache *entity.SongHashCache) bool {
	if data.Sha256 != "" {
		if md5, ok := cache.GetMD5(data.Sha256); ok {
			data.Md5 = md5
			return true
		}
		return false
	} else if data.Md5 != "" {
		if sha256, ok := cache.GetSHA256(data.Md5); ok {
			data.Sha256 = sha256
			return true
		}
		return false
	}
	return false
}

// For exporting usage (i.e imported by beatoraja)
//
// Nothing particular, but need to modify the output for beatoraja
type DiffTableDataExportDto struct {
	Level  string `json:"level"`
	Sha256 string `json:"sha256"`
	Md5    string `json:"md5"`
	Title  string `json:"title"`
}
