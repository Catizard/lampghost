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
	db.Table("rival_score_log").AutoMigrate(&entity.RivalScoreLog{})
	db.Table("difftable_header").AutoMigrate(&entity.DiffTableHeader{})
	db.Table("difftable_data").AutoMigrate(&entity.DiffTableData{})
	db.Table("course_info").AutoMigrate(&entity.CourseInfo{})
	db.Table("rival_song_data").AutoMigrate(&entity.RivalSongData{})
	db.Table("rival_tag").AutoMigrate(&entity.RivalTag{})

	log.Debugf("Initialized database at %s\n", config.DSN)
	return db, err
}
