package difftable

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Catizard/lampghost/internel/common"
	"github.com/Catizard/lampghost/internel/tui/choose"
	_ "github.com/mattn/go-sqlite3"
)

type DiffTableHeader struct {
	Id           int    `db:"id"`
	DataUrl      string `json:"data_url" db:"data_url"`
	DataLocation string `db:"data_location"`
	LastUpdate   string `json:"last_update" db:"last_update"`
	Name         string `json:"name" db:"name"`
	OriginalUrl  string `json:"original_url" db:"original_url"`
	Symbol       string `json:"symbol" db:"symbol"`
	Alias        string `json:"alias" db:"alias"`
}

func (header *DiffTableHeader) String() string {
	if len(header.Alias) > 0 {
		return fmt.Sprintf("%s(%s) [symbol=%s, url=%s]", header.Name, header.Alias, header.Symbol, header.Alias)
	}
	return fmt.Sprintf("%s [symbol=%s, url=%s]", header.Name, header.Symbol, header.Alias)
}

// Initialize difftable_header table
func InitDifftableHeaderTable() error {
	db := common.OpenDB()
	defer db.Close()
	_, err := db.Exec("DROP TABLE IF EXISTS 'difftable_header';CREATE TABLE difftable_header ( id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, alias TEXT, last_update TEXT, symbol TEXT NOT NULL, data_location TEXT NOT NULL);")
	return err
}

// Add a difficult table header(or say, meta data) and related song data to disk.
func (header *DiffTableHeader) AddDiffTable() error {
	// 1. Validation
	if arr, err := QueryDifficultTableHeaderByName(header.Name); err != nil || len(arr) > 0 {
		if err != nil {
			return err
		}
		log.Fatalf(`There is already a table named (or its alias matches) %s
		Use table sync command to sync table.
		Use table del command to delete table`, header.Name)
	}
	// 2. Create data file
	if err := saveTableData(header.getDataJsonFileName(), header.DataUrl); err != nil {
		return err
	}
	// 3. Insert into database
	dataLocation := header.getDataJsonFileName()
	header.DataLocation = dataLocation
	db := common.OpenDB()
	defer db.Close()
	if _, err := db.NamedExec(`INSERT INTO difftable_header(id, data_url, data_location, last_update, name, symbol, alias) VALUES (:id, :data_url, :data_location, :last_update, :name, :symbol, :alias)`, header); err != nil {
		return err
	}
	return nil
}

// Query by name or alias
func QueryDifficultTableHeaderByName(name string) ([]DiffTableHeader, error) {
	db := common.OpenDB()
	defer db.Close()
	var ret []DiffTableHeader
	err := db.Select(&ret, "SELECT * FROM difftable_header WHERE name=? OR alias=?", name, name)
	return ret, err
}

// Simple choose wrapper of QueryDifficultTableHeaderByName
func QueryDifficultTableHeaderByNameWithChoices(name string) (DiffTableHeader, error) {
	dthArr, err := QueryDifficultTableHeaderByName(name)
	if err != nil {
		return DiffTableHeader{}, err
	}
	choices := make([]string, 0)
	for _, v := range dthArr {
		choices = append(choices, v.String())
	}
	i := choose.OpenChooseTui(choices, "Multiple table matched with %s, please choose one:")
	return dthArr[i], nil
}

func (header *DiffTableHeader) SyncDifficultTable() error {
	// TODO: header's info might changed, sync header part also

	// sync data.json, propagate any error
	return saveTableData(header.getDataJsonFileName(), header.DataUrl)
}

// Download specified difficult table's data to disk
// File would be overwrite if it already exist
func saveTableData(fileName string, dataUrl string) error {
	// 1. Create data file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// 2. Download to file
	// TODO: if dataUrl is not start with http...
	resp, err := http.Get(dataUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(file, resp.Body)
	return nil
}

// Return difficult table header's data json file name
func (header *DiffTableHeader) getDataJsonFileName() string {
	return header.Name + ".json"
}
