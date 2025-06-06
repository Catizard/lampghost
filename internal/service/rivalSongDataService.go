package service

import (
	"sync"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

var (
	defaultSongHashCache     *entity.SongHashCache = nil
	defaultSongHashCacheLock sync.Mutex            = sync.Mutex{}
)

type RivalSongDataService struct {
	db *gorm.DB
}

func NewRivalSongDataService(db *gorm.DB) *RivalSongDataService {
	return &RivalSongDataService{
		db: db,
	}
}

// Basic query function for rival_song_data table
func findRivalSongDataByID(tx *gorm.DB, ID uint) (*entity.RivalSongData, error) {
	var data *entity.RivalSongData
	if err := tx.Find(&data, ID).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Basic query function for rival_song_data table
func findRivalSongDataList(tx *gorm.DB, filter *entity.RivalSongData) ([]*entity.RivalSongData, int, error) {
	var songDataList []*entity.RivalSongData
	partial := tx.Model(&entity.RivalSongData{})
	if filter != nil {
		partial = partial.Where(filter)
	}
	if err := partial.Find(&songDataList).Error; err != nil {
		return nil, 0, eris.Wrap(err, "cannot select from rival_song_data")
	}
	return songDataList, len(songDataList), nil
}

// Query the default song hash cache, which is built by main user's save file
func queryDefaultSongHashCache(tx *gorm.DB) (*entity.SongHashCache, error) {
	defaultSongHashCacheLock.Lock()
	defer defaultSongHashCacheLock.Unlock()
	// TODO: if cache build was failed, then it would repeat itself until it success
	// it might be a performance problem since most cases you could never build the cache
	// if it fails the first time
	if defaultSongHashCache != nil {
		return defaultSongHashCache, nil
	}
	mainUser, err := queryMainUser(tx)
	if err != nil {
		return nil, eris.Wrap(err, "failed to query main user")
	}
	cache, err := querySongHashCache(tx, mainUser.ID)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to query song hash cache from user(id=%d)", mainUser.ID)
	}
	defaultSongHashCache = cache
	return cache, nil
}

// Should be called every time there is a change in the data of songdata table
func expireDefaultCache() {
	defaultSongHashCacheLock.Lock()
	defer defaultSongHashCacheLock.Unlock()
	defaultSongHashCache = nil
	log.Debugf("[RivalSongDataService] Expired default song cache")
}

// Build song hash cache by specified user's save file
func querySongHashCache(tx *gorm.DB, rivalID uint) (*entity.SongHashCache, error) {
	md5KeyCache := make(map[string]string)
	sha256KeyCache := make(map[string]string)
	dataList, _, err := findRivalSongDataList(tx, &entity.RivalSongData{RivalId: rivalID})
	if err != nil {
		return nil, eris.Wrapf(err, "failed to query user(id=%d)'s songdata.db contents", rivalID)
	}
	for _, data := range dataList {
		md5KeyCache[data.Md5] = data.Sha256
		sha256KeyCache[data.Sha256] = data.Md5
	}
	return entity.NewSongHashCache(md5KeyCache, sha256KeyCache), nil
}

// Build song hash cache by the contents which is passed as a parameter
//
// Used when rival_song_data table is not ready, e.g. initialize phase
func generateSongHashCacheFromRawData(songData []*entity.SongData) *entity.SongHashCache {
	md5KeyCache := make(map[string]string)
	sha256KeyCache := make(map[string]string)
	for _, data := range songData {
		md5KeyCache[data.Md5] = data.Sha256
		sha256KeyCache[data.Sha256] = data.Md5
	}
	return entity.NewSongHashCache(md5KeyCache, sha256KeyCache)
}
