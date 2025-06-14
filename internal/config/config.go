package config

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
)

const (
	DBFileName = "lampghost.db"
)

var (
	WorkingDirectory string
	lock             sync.Mutex
)

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
	viper.SetDefault("IgnoreVariantCourse", 0)
	viper.SetDefault("Locale", "en")
	viper.SetDefault("DownloadSite", "wriggle")
	viper.SetDefault("SeparateDownloadMD5", "https://bms.wrigglebug.xyz/download/package/%s")
	viper.SetDefault("MaximumDownloadCount", 5)
	viper.SafeWriteConfig()
}

type ApplicationConfig struct {
	InternalServerPort  int32
	IgnoreVariantCourse int32
	Locale              string
	// Constants, currently the only option is wriggle
	DownloadSite string
	// Related to DownloadSite
	SeparateDownloadMD5 string
	// Download directory
	DownloadDirectory    string
	MaximumDownloadCount int
	// Extra fields, need to be setted explicitly
	EnableDownloadFeature bool
}

func ReadConfig() (*ApplicationConfig, error) {
	lock.Lock()
	defer lock.Unlock()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	conf := &ApplicationConfig{
		InternalServerPort:   viper.GetInt32("InternalServerPort"),
		IgnoreVariantCourse:  viper.GetInt32("IgnoreVariantCourse"),
		Locale:               viper.GetString("Locale"),
		DownloadSite:         viper.GetString("DownloadSite"),
		SeparateDownloadMD5:  viper.GetString("SeparateDownloadMD5"),
		DownloadDirectory:    viper.GetString("DownloadDirectory"),
		MaximumDownloadCount: viper.GetInt("MaximumDownloadCount"),
	}
	conf.EnableDownloadFeature = conf.EnableDownload() == nil
	return conf, nil
}

func (c *ApplicationConfig) WriteConfig() error {
	lock.Lock()
	defer lock.Unlock()
	viper.Set("InternalServerPort", c.InternalServerPort)
	viper.Set("IgnoreVariantCourse", c.IgnoreVariantCourse)
	viper.Set("Locale", c.Locale)
	viper.Set("DownloadSite", c.DownloadSite)
	viper.Set("SeparateDownloadMD5", c.SeparateDownloadMD5)
	viper.Set("DownloadDirectory", c.DownloadDirectory)
	viper.Set("MaximumDownloadCount", c.MaximumDownloadCount)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

// Query if download related config options are well setted
func (c *ApplicationConfig) EnableDownload() error {
	if c.DownloadDirectory == "" {
		return eris.New("cannot download: download directory hasn't been set yet")
	}
	if c.SeparateDownloadMD5 == "" {
		return eris.New("cannot download: download url hasn't been set yet")
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
