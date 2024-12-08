package database

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	db  *gorm.DB
	DSN string
}

func NewDatabase(DSN string) *DB {
	return &DB{
		DSN: DSN,
	}
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

func (db *DB) Automigrate() {
	db.db.Table("rival_info").AutoMigrate(&entity.RivalInfo{})
}
