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
		if err := conf.WriteConfig(); err != nil {
			return err
		}
		return nil
	})
}
