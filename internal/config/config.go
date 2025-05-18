package config

import (
	"errors"
	"os"
	"path/filepath"

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
	// Create the directory if it's not exist
	_, err = os.Stat(WorkingDirectory)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(WorkingDirectory), 0700); err != nil {
			panic(err)
		}
	}
	// Setup viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(WorkingDirectory[:len(WorkingDirectory)-1])
	// Setup defaults
	viper.SetDefault("InternalServerPort", 7391)
	viper.SetDefault("FolderSymbol", "")
	viper.SetDefault("IgnoreVariantCourse", 0)
	viper.SetDefault("Locale", "en")
	viper.SafeWriteConfig()
}

type ApplicationConfig struct {
	InternalServerPort  int32
	FolderSymbol        string
	IgnoreVariantCourse int32
	Locale              string
}

func ReadConfig() (*ApplicationConfig, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	return &ApplicationConfig{
		InternalServerPort:  viper.GetInt32("InternalServerPort"),
		FolderSymbol:        viper.GetString("FolderSymbol"),
		IgnoreVariantCourse: viper.GetInt32("IgnoreVariantCourse"),
		Locale:              viper.GetString("Locale"),
	}, nil
}

func (c *ApplicationConfig) WriteConfig() error {
	viper.Set("InternalServerPort", c.InternalServerPort)
	viper.Set("FolderSymbol", c.FolderSymbol)
	viper.Set("IgnoreVariantCourse", c.IgnoreVariantCourse)
	viper.Set("Locale", c.Locale)
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
