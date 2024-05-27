package common

import (
	"errors"
	"log"
	"os"
)

const (
	DBFileName = "lampghost.db"
)

// Panic if lampghost.db is not exist
func CheckInitialize() {
	if _, err := os.Stat(DBFileName); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Call init command before you do anything")
	}
}
