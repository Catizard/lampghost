package database

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	db  *gorm.DB
	DSN string
}

func NewDatabase(config *config.DatabaseConfig) (*DB, error) {
	out := &DB{
		DSN: config.DSN,
	}

	if err := out.Open(); err != nil {
		return nil, err
	}

	if err := out.Automigrate(); err != nil {
		return nil, err
	}

	log.Debugf("Initialized database at %s\n", out.DSN)
	return out, nil
}

func (db *DB) Open() (err error) {
	if db.DSN == "" {
		return fmt.Errorf("DSN cannot be empty")
	}

	if db.db, err = gorm.Open(sqlite.Open(db.DSN), &gorm.Config{}); err != nil {
		return err
	}

	return err
}

func (db *DB) Automigrate() error {
	return db.db.Table("rival_info").AutoMigrate(&entity.RivalInfo{})
}
