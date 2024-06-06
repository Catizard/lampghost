package common

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DBFileName = "lampghost.db"
)

// Open lampghost database
// Panic if any error occurred
func OpenDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", DBFileName)
	if err != nil {
		panic(err)
	}
	return db
}
