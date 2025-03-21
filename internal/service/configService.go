package service

import (
	"github.com/Catizard/lampghost_wails/internal/config"
	"gorm.io/gorm"
)

type ConfigService struct {
	db *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{
		db: db,
	}
}

func (s *ConfigService) WriteConfig(conf *config.ApplicationConfig) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		prevConf, err := config.ReadConfig()
		if err != nil {
			return err
		}
		if shouldReload(prevConf, conf) {
			mainUser, err := queryMainUser(tx)
			if err != nil {
				return err
			}
			mainUser.ScoreLogPath = &conf.ScorelogFilePath
			mainUser.SongDataPath = &conf.SongdataFilePath
			if err := updateRivalInfo(tx, mainUser); err != nil {
				return err
			}
			// TODO: mainUser.scorePath = &conf.ScoreFilePath
			if err := syncRivalData(tx, mainUser); err != nil {
				return err
			}
		}
		if err := conf.WriteConfig(); err != nil {
			return err
		}
		return nil
	})
}

func shouldReload(prevConf *config.ApplicationConfig, newConf *config.ApplicationConfig) bool {
	if prevConf.ScorelogFilePath != newConf.ScorelogFilePath {
		return true
	}
	if prevConf.SongdataFilePath != newConf.SongdataFilePath {
		return true
	}
	if prevConf.ScoreFilePath != newConf.ScoreFilePath {
		return true
	}
	return false
}
