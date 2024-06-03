package common

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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

// Open lampghost database
// Panic if any error occurred
func OpenDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", DBFileName)
	if err != nil {
		panic(err)
	}
	return db
}
