package service

import (
	"sync"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

var defaultSongHashCache *entity.SongHashCache = nil
var defaultSongHashCacheLock sync.Mutex = sync.Mutex{}

type RivalSongDataService struct {
	db *gorm.DB
}

func NewRivalSongDataService(db *gorm.DB) *RivalSongDataService {
	return &RivalSongDataService{
		db: db,
	}
}

func (s *RivalSongDataService) BindRivalSongDataToFolder(rivalSongDataID uint, folderIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		songData, err := findRivalSongDataByID(tx, rivalSongDataID)
		if err != nil {
			return err
		}

		content := entity.FolderContent{
			Sha256: songData.Sha256,
			Md5:    songData.Md5,
			Title:  songData.Title,
		}

		if err := bindSongToFolder(tx, content, folderIDs); err != nil {
			return err
		}

		return nil
	})
}

func findRivalSongDataList(tx *gorm.DB, filter *entity.RivalSongData) ([]entity.RivalSongData, int, error) {
	var songDataList []entity.RivalSongData
	if err := tx.Where(filter).Find(&songDataList).Error; err != nil {
		return nil, 0, err
	}
	return songDataList, len(songDataList), nil
}

func findRivalSongDataByID(tx *gorm.DB, ID uint) (*entity.RivalSongData, error) {
	var data *entity.RivalSongData
	if err := tx.Find(&data, ID).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func queryRivalSongDataExistence(tx *gorm.DB, filter *entity.RivalSongData) (bool, error) {
	var dupCount int64
	if err := tx.Model(&entity.RivalSongData{}).Where(filter).Count(&dupCount).Error; err != nil {
		return false, err
	}
	return dupCount > 0, nil
}

func queryDefaultSongHashCache(tx *gorm.DB) (*entity.SongHashCache, error) {
	defaultSongHashCacheLock.Lock()
	defer defaultSongHashCacheLock.Unlock()
	// TODO: if cache build was failed, then it would repeat itself until it success
	// it might be a performance problem since most cases you could never build the cache
	if defaultSongHashCache != nil {
		return defaultSongHashCache, nil
	}
	mainUser, err := queryMainUser(tx)
	if err != nil {
		return nil, err
	}
	cache, err := querySongHashCache(tx, mainUser.ID)
	if err != nil {
		return nil, err
	}
	defaultSongHashCache = cache
	return cache, nil
}

// Should be called every time there is a change in the data of songdata table
func expireDefaultCache() {
	defaultSongHashCacheLock.Lock()
	defer defaultSongHashCacheLock.Unlock()
	defaultSongHashCache = nil
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
