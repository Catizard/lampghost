package config

import (
	"os"

	"github.com/Catizard/lampghost/internal/common"
	"github.com/charmbracelet/log"
)

var WorkingDirectory string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Make it configurable
	WorkingDirectory = homeDir + "/.lampghost/"
}

// Get LampGhost application's datasource name
// LampGhost will keep using one single datasource, so this is fine
func GetDSN() string {
	return WorkingDirectory + common.DBFileName
}
