package dto

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
)

type DiffTableDataDto struct {
	ID       uint
	HeaderID uint
	Artist   string
	Comment  string
	Level    string
	Lr2BmsId string
	Md5      string
	NameDiff string
	Title    string
	Url      string
	UrlDiff  string
	Sha256   string

	Lamp      int
	PlayCount int

	GhostLamp      int
	GhostPlayCount int
	DataLost       bool
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
	ok := ret.RepairHash(cache)
	ret.DataLost = !ok
	return ret
}

// Repair the data that has only sha256 or md5
//
// Return whether this data could be found in cache or not. In other words, do we have this song?
// The reason that 'repair' step and 'check' step are implemented seperately is to ensure correctness
func (data *DiffTableDataDto) RepairHash(cache *entity.SongHashCache) bool {
	if data.Sha256 != "" {
		if md5, ok := cache.GetMD5(data.Sha256); ok {
			data.Md5 = md5
		}
	} else if data.Md5 != "" {
		if sha256, ok := cache.GetSHA256(data.Md5); ok {
			data.Sha256 = sha256
		}
	}

	if data.Sha256 != "" {
		if _, ok := cache.GetMD5(data.Sha256); ok {
			return true
		}
	}
	if data.Md5 != "" {
		if _, ok := cache.GetSHA256(data.Md5); ok {
			return true
		}
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
