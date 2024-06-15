package config

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
)

const (
	DBFileName = "lampghost.db"
)

var WorkingDirectory string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Make it configurable
	// Note: Almost every data should be put in working directory
	// The only exception is rival's log files, which follows with executable file
	WorkingDirectory = homeDir + "/.lampghost/"
}

func CheckInitialize() {
	if _, err := os.Stat(GetDSN()); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Call init command before you do anything")
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
