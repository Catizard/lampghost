package config

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/Catizard/lampghost_wails/internal/config/download"
	"github.com/charmbracelet/log"
	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
)

const (
	DBFileName = "lampghost.db"
	VERSION    = "0.2.6"
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
	viper.SetDefault("DownloadSite", download.DefaultDownloadSource.GetMeta().Name)
	viper.SetDefault("PreviewSite", DefaultPreviewSource.GetName())
	// viper.SetDefault("SeparateDownloadMD5", "https://bms.wrigglebug.xyz/download/package/%s")
	viper.SetDefault("MaximumDownloadCount", 5)
	viper.SetDefault("EnableAutoReload", 1)
	viper.SafeWriteConfig()
	// Setup logger
	if logFile, err := os.OpenFile(WorkingDirectory+"lampghost.log", os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		panic(err)
	} else {
		log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	}
}

type ApplicationConfig struct {
	InternalServerPort  int32
	IgnoreVariantCourse int32
	Locale              string
	// Download source: wriggle | konmai
	DownloadSite string
	// Related to DownloadSite
	// SeparateDownloadMD5 string
	// Download directory
	DownloadDirectory    string
	MaximumDownloadCount int
	// Auto reload save files when having write operation on scorelog.db
	EnableAutoReload int
	// Preview source: sayaka | konmai
	PreviewSite string
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
		InternalServerPort:  viper.GetInt32("InternalServerPort"),
		IgnoreVariantCourse: viper.GetInt32("IgnoreVariantCourse"),
		Locale:              viper.GetString("Locale"),
		DownloadSite:        viper.GetString("DownloadSite"),
		// SeparateDownloadMD5:  viper.GetString("SeparateDownloadMD5"),
		DownloadDirectory:    viper.GetString("DownloadDirectory"),
		MaximumDownloadCount: viper.GetInt("MaximumDownloadCount"),
		EnableAutoReload:     viper.GetInt("EnableAutoReload"),
		PreviewSite:          viper.GetString("PreviewSite"),
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
	// viper.Set("SeparateDownloadMD5", c.SeparateDownloadMD5)
	viper.Set("DownloadDirectory", c.DownloadDirectory)
	viper.Set("MaximumDownloadCount", c.MaximumDownloadCount)
	viper.Set("EnableAutoReload", c.EnableAutoReload)
	viper.Set("PreviewSite", c.PreviewSite)
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
	// if c.SeparateDownloadMD5 == "" {
	// 	return eris.New("cannot download: download url hasn't been set yet")
	// }
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

// Get LampGhost application's datasource name
// LampGhost will keep using one single datasource, so this is fine
func GetDSN() string {
	return WorkingDirectory + DBFileName
}

func JoinWorkingDirectory(relativePath string) string {
	return WorkingDirectory + relativePath
}
