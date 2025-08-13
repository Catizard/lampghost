package database

import (
	"fmt"
	"os"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	if err := db.Table("custom_diff_table").AutoMigrate(&entity.CustomDiffTable{}); err != nil {
		return err
	}

	if err := db.Table("custom_course").AutoMigrate(&entity.CustomCourse{}); err != nil {
		return err
	}

	if err := db.Table("custom_course_data").AutoMigrate(&entity.CustomCourseData{}); err != nil {
		return err
	}

	if err := db.Table("song_directory").AutoMigrate(&entity.SongDirectory{}); err != nil {
		return err
	}

	if err := db.Table("rival_score_data").AutoMigrate(&entity.RivalScoreData{}); err != nil {
		return err
	}

	// I cannot find a better solution
	var defaultCustomTable entity.CustomDiffTable
	defaultCustomTable.ID = 1
	if err := db.FirstOrCreate(&defaultCustomTable, entity.CustomDiffTable{
		Model: gorm.Model{
			ID: 1,
		},
		Name:        "lampghost",
		Symbol:      "",
		LevelOrders: "",
	}).Error; err != nil {
		panic(err)
	}

	return nil
}

func NewDatabase(config *config.DatabaseConfig) (db *gorm.DB, err error) {
	newLogger := logger.New(log.Default(), logger.Config{})
	if db, err = gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
		Logger: newLogger,
	}); err != nil {
		return nil, err
	}
	if err := migrates(db); err != nil {
		return nil, err
	}
	log.Debugf("Initialized database at %s\n", config.DSN)
	return db, err
}

// Open the connection to self generated songdata.db file
func NewSelfGeneratedSongDataDatabase(purge bool) (db *gorm.DB, err error) {
	fp := config.GetSelfGeneratedSongDataPath()
	if purge {
		os.Remove(fp)
	}
	if db, err = gorm.Open(sqlite.Open(fp)); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&entity.SongData{}); err != nil {
		return nil, err
	}
	return
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
