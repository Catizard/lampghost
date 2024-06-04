package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB represents database connection
type DB struct {
	db *sqlx.DB
	// data source name
	DSN string
}

func NewDB(dsn string) *DB {
	db := &DB{
		DSN: dsn,
	}
	return db
}

func (db *DB) Open() (err error) {
	if db.DSN == "" {
		return fmt.Errorf("panic: sqlite::DB.open()")
	}

	if db.db, err = sqlx.Open("sqlite3", db.DSN); err != nil {
		return err
	}
	return nil
}

// Closes the database connection
func (db *DB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

// Wrapper of Sqlx.Tx object
type Tx struct {
	*sqlx.Tx
	db *DB
}