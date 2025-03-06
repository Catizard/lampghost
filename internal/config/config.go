package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

const (
	DBFileName = "lampghost.db"
)

var WorkingDirectory string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// TODO: Make it configurable
	WorkingDirectory = homeDir + "/.lampghost_wails/"
	// Setup viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	log.Infof("config directory: %s", WorkingDirectory[:len(WorkingDirectory)-1])
	viper.AddConfigPath(WorkingDirectory[:len(WorkingDirectory)-1])
	// Setup defaults
	viper.SetDefault("InternalServerPort", 7391)
	viper.SetDefault("FolderSymbol", "")
	viper.SetDefault("IgnoreVariantCourse", 0)
	viper.SetDefault("Locale", "en")
	viper.SetDefault("ForceFullyReload", 0)
	viper.SafeWriteConfig()
	// Create the directory if it's not exist
	_, err = os.Stat(WorkingDirectory)
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(WorkingDirectory), 0700); err != nil {
			panic(err)
		}
	}
}

type ApplicationConfig struct {
	UserName            string
	ScorelogFilePath    string
	SongdataFilePath    string
	ScoreFilePath       string
	InternalServerPort  int32
	FolderSymbol        string
	IgnoreVariantCourse int32
	Locale              string
	ForceFullyReload    int32
}

func ReadConfig() (*ApplicationConfig, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	return &ApplicationConfig{
		UserName:            viper.GetString("UserName"),
		ScorelogFilePath:    viper.GetString("ScoreLogFilePath"),
		SongdataFilePath:    viper.GetString("SongDataFilePath"),
		ScoreFilePath:       viper.GetString("ScoreFilePath"),
		InternalServerPort:  viper.GetInt32("InternalServerPort"),
		FolderSymbol:        viper.GetString("FolderSymbol"),
		IgnoreVariantCourse: viper.GetInt32("IgnoreVariantCourse"),
		Locale:              viper.GetString("Locale"),
		ForceFullyReload:    viper.GetInt32("ForceFullyReload"),
	}, nil
}

func (c *ApplicationConfig) WriteConfig() error {
	viper.Set("UserName", c.UserName)
	viper.Set("ScoreLogFilePath", c.ScorelogFilePath)
	viper.Set("SongDataFilePath", c.SongdataFilePath)
	viper.Set("ScoreFilePath", c.ScoreFilePath)
	viper.Set("InternalServerPort", c.InternalServerPort)
	viper.Set("FolderSymbol", c.FolderSymbol)
	viper.Set("IgnoreVariantCourse", c.IgnoreVariantCourse)
	viper.Set("Locale", c.Locale)
	viper.Set("ForceFullyReload", c.ForceFullyReload)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

type DatabaseConfig struct {
	DSN string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		DSN: WorkingDirectory + DBFileName,
	}
}

func CheckInitialize() {
	if _, err := os.Stat(GetDSN()); errors.Is(err, os.ErrNotExist) {
		panic("Call init command before you do anything")
	}
}

// Get LampGhost application's datasource name
// LampGhost will keep using one single datasource, so this is fine
func GetDSN() string {
	return WorkingDirectory + DBFileName
}

func JoinWorkingDirectory(relativePath string) string {
	return WorkingDirectory + relativePath
}
