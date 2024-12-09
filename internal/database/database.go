package database

import (
	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(config *config.DatabaseConfig) (db *gorm.DB, err error) {
	if db, err = gorm.Open(sqlite.Open(config.DSN), &gorm.Config{}); err != nil {
		return nil, err
	}

	db.Table("rival_info").AutoMigrate(&entity.RivalInfo{})

	log.Debugf("Initialized database at %s\n", config.DSN)
	return db, err
}
