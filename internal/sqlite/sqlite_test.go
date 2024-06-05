package sqlite_test

import (
	"testing"

	"github.com/Catizard/lampghost/internal/config"
	"github.com/Catizard/lampghost/internal/sqlite"
	"github.com/charmbracelet/log"
)

func TestInitialize(t *testing.T) {
	sqlite.InitializeDatabase()
	db := mustOpen()
	if n, err := sqlite.QueryTableCount(config.GetDSN()); err != nil {
		t.Fatal(err)
	} else if n != 5 {
		t.Fatal("expect: 5, get: ", n)
	}
	mustClose(db)
}

// Open a database or fail
func mustOpen() *sqlite.DB {
	db := sqlite.NewDB(config.GetDSN())
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	return db
}

// Close a database or fail
func mustClose(db *sqlite.DB) {
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}
