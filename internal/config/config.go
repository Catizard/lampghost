package config

import (
	"errors"
	"os"
	"path/filepath"
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
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(WorkingDirectory), 0700); err != nil {
			panic(err)
		}
	}
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
