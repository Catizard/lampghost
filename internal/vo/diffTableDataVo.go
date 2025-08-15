package vo

import (
	"time"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
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
	SortBy     *string
	// NOTE: NEVER access this field directly, use GetOrder() instead
	SortOrder *string
	// Extra filter fields
	// See config.go for details
	UseScoredataForMainUser bool
	IDs                     []uint
	HeaderIDs               []uint

	RivalID      uint
	GhostRivalID uint
	// See diffTableDataService#findDiffTableDataListWithRival for usage
	EndGhostRecordTime time.Time
	Md5s               []string
	Levels             []string
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

func (data *DiffTableDataVo) GetOrder() string {
	if data.SortOrder == nil {
		return "asc"
	}
	switch *data.SortOrder {
	case "ascend":
		return "asc"
	case "descend":
		return "desc"
	default:
		log.Warnf("unexpected SortOrder value: %s", *data.SortOrder)
		return "asc"
	}
}
