package database

import (
	"fmt"
	"os"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// database auto migrate definition
func migrates(db *gorm.DB) error {
	if err := db.Table("rival_info").AutoMigrate(&entity.RivalInfo{}); err != nil {
		return err
	}

	if err := db.Table("rival_score_log").AutoMigrate(&entity.RivalScoreLog{}); err != nil {
		return err
	}

	if err := db.Table("difftable_header").AutoMigrate(&entity.DiffTableHeader{}); err != nil {
		return err
	}

	if err := db.Table("difftable_data").AutoMigrate(&entity.DiffTableData{}); err != nil {
		return err
	}

	if err := db.Table("course_info").AutoMigrate(&entity.CourseInfo{}); err != nil {
		return err
	}

	if err := db.Table("rival_song_data").AutoMigrate(&entity.RivalSongData{}); err != nil {
		return err
	}

	if err := db.Table("rival_tag").AutoMigrate(&entity.RivalTag{}); err != nil {
		return err
	}

	if err := db.Table("folder").AutoMigrate(&entity.Folder{}); err != nil {
		return err
	}

	if err := db.Table("folder_content").AutoMigrate(&entity.FolderContent{}); err != nil {
		return err
	}

	if err := db.Table("rival_score_data_log").AutoMigrate(&entity.RivalScoreDataLog{}); err != nil {
		return err
	}

	if err := db.Table("predefine_table_header").AutoMigrate(&entity.PredefineTableHeader{}); err != nil {
		return err
	}

	return nil
}

func NewDatabase(config *config.DatabaseConfig) (db *gorm.DB, err error) {
	if db, err = gorm.Open(sqlite.Open(config.DSN), &gorm.Config{}); err != nil {
		return nil, err
	}
	if err := migrates(db); err != nil {
		return nil, err
	}
	log.Debugf("Initialized database at %s\n", config.DSN)
	return db, err
}

// Create a memory database for tests
func NewMemoryDatabase() (db *gorm.DB, err error) {
	if db, err = gorm.Open(sqlite.Open(":memory:")); err != nil {
		return nil, err
	}
	if err := migrates(db); err != nil {
		return nil, err
	}
	return db, err
}

// Helper function for validating local databasement file path
func VerifyLocalDatabaseFilePath(filePath string) error {
	if stat, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("assert: no file exists at %s", filePath)
		}
		return fmt.Errorf("assert: cannot stat file at %s", filePath)
	} else if stat.IsDir() {
		return fmt.Errorf("assert: file path %s is a directory, not an valid database file", filePath)
	}
	return nil
}
