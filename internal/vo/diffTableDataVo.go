package vo

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type DiffTableDataVo struct {
	gorm.Model
	HeaderID uint
	Artist   string
	Comment  string
	Level    string
	Lr2BmsId string `json:"lr2_bmdid"`
	Md5      string
	NameDiff string
	Title    string
	Url      string `json:"url"`
	UrlDiff  string `json:"url_diff"`
	Sha256   string

	Pagination *entity.Page
	// Extra filter fields
	IDs       []uint
	HeaderIDs []uint
}

func (data *DiffTableDataVo) Entity() *entity.DiffTableData {
	return &entity.DiffTableData{
		Model: gorm.Model{
			ID:        data.ID,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			DeletedAt: data.DeletedAt,
		},
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
