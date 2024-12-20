package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

type RivalSongDataService struct {
	db                   *gorm.DB
	defaultSongHashCache *SongHashCache
}

func NewRivalSongDataService(db *gorm.DB, rivalInfoSerivce *RivalInfoService) *RivalSongDataService {
	return &RivalSongDataService{
		db: db,
	}
}

func (s *RivalSongDataService) FindRivalSongDataList(filter *entity.RivalSongData) ([]entity.RivalSongData, int, error) {
	var songDataList []entity.RivalSongData
	if err := s.db.Where(filter).Find(&songDataList).Error; err != nil {
		return nil, 0, err
	}
	return songDataList, len(songDataList), nil
}

func (s *RivalSongDataService) QuerySongHashCache(rivalId uint) (*SongHashCache, error) {
	md5KeyCache := make(map[string]string)
	sha256KeyCache := make(map[string]string)
	dataList, _, err := s.FindRivalSongDataList(&entity.RivalSongData{RivalId: rivalId})
	if err != nil {
		return nil, err
	}
	for _, data := range dataList {
		md5KeyCache[data.Md5] = data.Sha256
		sha256KeyCache[data.Sha256] = data.Md5
	}
	return &SongHashCache{
		md5KeyCache:    md5KeyCache,
		sha256KeyCache: sha256KeyCache,
	}, nil
}

type SongHashCache struct {
	md5KeyCache    map[string]string
	sha256KeyCache map[string]string
}

func (c *SongHashCache) GetSHA256(md5 string) (string, bool) {
	v, ok := c.md5KeyCache[md5]
	return v, ok
}

func (c *SongHashCache) GetMD5(sha256 string) (string, bool) {
	v, ok := c.sha256KeyCache[sha256]
	return v, ok
}
