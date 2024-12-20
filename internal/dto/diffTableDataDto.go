package dto

import "github.com/Catizard/lampghost_wails/internal/entity"

type DiffTableDataDto struct {
	ID        uint
	HeaderID  uint
	Artist    string
	Comment   string
	Level     string
	Lr2BmsId  string `json:"lr2_bmdid"`
	Md5       string
	NameDiff  string
	Title     string
	Url       string `json:"url"`
	UrlDiff   string `json:"url_diff"`
	Sha256    string
	Lamp      int
	GhostLamp int
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
