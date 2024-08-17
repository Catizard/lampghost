package sqlite

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/Catizard/lampghost/internal/config"
	"github.com/Catizard/lampghost/internal/data/filter"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

//go:embed migration/*.sql
var migrationFS embed.FS

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

	// Note: One migrate implementation is to call migrate every time when open database
	// But I don't think this is a good option, so I choose another way:
	// Force the user calling init command to migrate database
	// if err := db.migrate(); err != nil {
	// 	return err
	// }
	return nil
}

// Initialize database and working directory for lampghost
// Executes every sql file under 'sqlite/migration/'
// This method insure all or nothing (expect directory)
func InitializeDatabase() error {
	os.RemoveAll(filepath.Dir(config.WorkingDirectory))
	// Create working directory
	if err := os.MkdirAll(filepath.Dir(config.WorkingDirectory), 0700); err != nil {
		return err
	}

	// Initialize database
	db := NewDB(config.GetDSN())
	if err := db.Open(); err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.db.Exec("CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);"); err != nil {
		return fmt.Errorf("panic: cannot create migrations table: %w", err)
	}

	names, err := fs.Glob(migrationFS, "migration/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)
	if len(names) != 1 {
		log.Fatal("panic: migration sql files count != 1")
	}

	// Open transaction
	tx, err := db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, name := range names {
		if err := migrateFile(tx, name); err != nil {
			return fmt.Errorf("migration error: name=%q, err=%w", name, err)
		}
	}
	return tx.Commit()
}

// Query table count in database
// return -1 if error
func QueryTableCount(dsn string) (int, error) {
	db := NewDB(dsn)
	if err := db.Open(); err != nil {
		return -1, err
	}
	defer db.Close()
	var n int
	if err := db.db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name != 'android_metadata' AND name != 'sqlite_sequence';").Scan(&n); err != nil {
		return -1, err
	}
	return n, nil
}

func migrateFile(tx *Tx, name string) error {
	var n int
	if err := tx.QueryRow("SELECT COUNT(*) FROM migrations WHERE name=?", name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil // skip
	}

	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	_, err := tx.Exec("INSERT INTO migrations(name) VALUES(?)", name)
	return err
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

func (db *DB) BeginTx() (*Tx, error) {
	tx, err := db.db.Beginx()
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx: tx,
		db: db,
	}, nil
}

// Boilerplate code that directly read rows from external database
func DirectlyLoadTable[T interface{}](path string, tb string) ([]*T, error) {
	return DirectlyLoadTableWithFilter[T](path, tb, filter.NullFilter)	
}

func DirectlyLoadTableWithFilter[T interface{}](path string, tb string, filter null.Value[filter.Filter]) ([]*T, error) {
	// database initialize
	db := NewDB(path)
	if err := db.Open(); err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	boilerplateSql := fmt.Sprintf("SELECT * FROM %s", tb)
	if filter.Valid {
		boilerplateSql += " " + filter.ValueOrZero().GenerateWhereClause()
	}
	log.Debugf("boilerplateSql=%s", boilerplateSql)
	rows, err := tx.Queryx(boilerplateSql)
	if err != nil {
		return nil, err
	}

	ret := make([]*T, 0)
	for rows.Next() {
		var obj T
		err = rows.StructScan(&obj)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &obj)
	}
	return ret, nil
}
