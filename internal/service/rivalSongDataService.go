package service

import (
	"sync"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
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

func (s *RivalSongDataService) QueryDefaultSongHashCache() (*entity.SongHashCache, error) {
	return queryDefaultSongHashCache(s.db)
}

func (s *RivalSongDataService) QuerySongDataPageList(filter *vo.RivalSongDataVo) (out []*dto.RivalSongDataDto, cnt int, err error) {
	err = s.db.Transaction(func(tx *gorm.DB) error {
		out, cnt, err = findRivalSongDataList(tx, filter)
		return err
	})
	return
}

func (s *RivalSongDataService) ReloadRivalSongData() error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return reloadRivalSongData(tx)
	})
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
func findRivalSongDataList(tx *gorm.DB, filter *vo.RivalSongDataVo) (out []*dto.RivalSongDataDto, cnt int, err error) {
	moved := tx.Select(`rival_song_data.*`).Model(&entity.RivalSongData{}).Scopes(scopeRivalSongDataFilter(filter))
	if filter != nil {
		moved = moved.Scopes(pagination(filter.Pagination))
	}
	if err = moved.Find(&out).Error; err != nil {
		return
	}
	cnt = len(out)
	return
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
	dataList, _, err := findRivalSongDataList(tx, &vo.RivalSongDataVo{RivalId: rivalID})
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

// fully reload rival_song_data
func reloadRivalSongData(tx *gorm.DB) error {
	mainUser, err := queryMainUser(tx)
	if err != nil {
		return err
	}
	fp := mainUser.SongDataPath
	rawSongData, err := loadSongData(*fp)
	if err != nil {
		return err
	}
	if err := syncSongData(tx, rawSongData, mainUser.ID); err != nil {
		return err
	}
	// invalidate default song cache since we have rebuilt the `rival_song_data` table
	expireDefaultCache()
	return nil
}

// Specialized scope for vo.RivalSongDataVo
func scopeRivalSongDataFilter(filter *vo.RivalSongDataVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		return moved
	}
}
