package entity

import (
	"github.com/Catizard/bmstable"
	"gorm.io/gorm"
)

type DiffTableData struct {
	gorm.Model
	HeaderID uint
	Artist   string
	Comment  string
	Level    string
	Lr2BmsId string `json:"lr2_bmdid"`
	Md5      string `gorm:"index"`
	NameDiff string
	Title    string
	Url      string `json:"url"`
	UrlDiff  string `json:"url_diff"`
	Sha256   string
}

func (DiffTableData) TableName() string {
	return "difftable_data"
}

// Convert bmstable's type definition into internal one
func NewDiffTableDataFromImport(importData *bmstable.DifficultTableData) *DiffTableData {
	return &DiffTableData{
		Artist:   importData.Artist,
		Comment:  importData.Comment,
		Level:    importData.Level,
		Lr2BmsId: importData.Lr2BmsID,
		Md5:      importData.Md5,
		NameDiff: importData.NameDiff,
		Title:    importData.Title,
		Url:      importData.URL,
		UrlDiff:  importData.URLDiff,
		Sha256:   importData.Sha256,
	}
}
