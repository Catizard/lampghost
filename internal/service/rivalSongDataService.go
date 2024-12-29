package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalSongDataService struct {
	db *gorm.DB
}

func NewRivalSongDataService(db *gorm.DB) *RivalSongDataService {
	return &RivalSongDataService{
		db: db,
	}
}
func findRivalSongDataList(tx *gorm.DB, filter *entity.RivalSongData) ([]entity.RivalSongData, int, error) {
	var songDataList []entity.RivalSongData
	if err := tx.Where(filter).Find(&songDataList).Error; err != nil {
		return nil, 0, err
	}
	return songDataList, len(songDataList), nil
}

func queryDefaultSongHashCache(tx *gorm.DB) (*entity.SongHashCache, error) {
	mainUser, err := queryMainUser(tx)
	if err != nil {
		return nil, err
	}
	return querySongHashCache(tx, mainUser.ID)
}

func querySongHashCache(tx *gorm.DB, rivalID uint) (*entity.SongHashCache, error) {
	md5KeyCache := make(map[string]string)
	sha256KeyCache := make(map[string]string)
	dataList, _, err := findRivalSongDataList(tx, &entity.RivalSongData{RivalId: rivalID})
	if err != nil {
		return nil, err
	}
	for _, data := range dataList {
		md5KeyCache[data.Md5] = data.Sha256
		sha256KeyCache[data.Sha256] = data.Md5
	}
	return entity.NewSongHashCache(md5KeyCache, sha256KeyCache), nil
}

func generateSongHashCacheFromRawData(songData []*entity.SongData) *entity.SongHashCache {
	md5KeyCache := make(map[string]string)
	sha256KeyCache := make(map[string]string)
	for _, data := range songData {
		md5KeyCache[data.Md5] = data.Sha256
		sha256KeyCache[data.Sha256] = data.Md5
	}
	return entity.NewSongHashCache(md5KeyCache, sha256KeyCache)
}
