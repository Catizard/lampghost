package entity

import "gorm.io/gorm"

type DiffTableData struct {
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
}

func (DiffTableData) TableName() string {
	return "difftable_data"
}

func (data *DiffTableData) RepairHash(cache *SongHashCache) bool {
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
