package service

import (
	"github.com/Catizard/lampghost_wails/internal/config"
	"gorm.io/gorm"
)

type ConfigService struct {
	db      *gorm.DB
	publish chan any
}

func NewConfigService(db *gorm.DB, configPublishChannel chan any) *ConfigService {
	return &ConfigService{
		db:      db,
		publish: configPublishChannel,
	}
}

func (s *ConfigService) WriteConfig(conf *config.ApplicationConfig) error {
	if err := conf.WriteConfig(); err != nil {
		return err
	}
	s.publish <- 1
	return nil
}
