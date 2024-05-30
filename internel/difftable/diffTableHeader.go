package difftable

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Catizard/lampghost/internel/common"
	"github.com/Catizard/lampghost/internel/tui/choose"
	_ "github.com/mattn/go-sqlite3"
)

type DiffTableHeader struct {
	Id           int            `db:"id"`
	DataUrl      string         `json:"data_url" db:"data_url"`
	DataLocation string         `db:"data_location"`
	LastUpdate   string         `json:"last_update" db:"last_update"`
	Name         string         `json:"name" db:"name"`
	OriginalUrl  string         `json:"original_url" db:"original_url"`
	Symbol       string         `json:"symbol" db:"symbol"`
	Alias        string         `json:"alias" db:"alias"`
	Course       [][]CourseInfo `json:"course"`
}

func (header *DiffTableHeader) String() string {
	if len(header.Alias) > 0 {
		return fmt.Sprintf("%s(%s) [symbol=%s, url=%s]", header.Name, header.Alias, header.Symbol, header.Alias)
	}
	return fmt.Sprintf("%s [symbol=%s, url=%s]", header.Name, header.Symbol, header.Alias)
}

// Initialize difftable_header table
func InitDiffTableHeaderTable() error {
	db := common.OpenDB()
	defer db.Close()
	_, err := db.Exec("DROP TABLE IF EXISTS 'difftable_header';CREATE TABLE difftable_header ( id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, alias TEXT, last_update TEXT, symbol TEXT NOT NULL, data_location TEXT NOT NULL, data_url TEXT NOT NULL);")
	return err
}

// Fetch difficult table header from url
// For now, it only supports .json file
func FetchDiffTableHeader(url string) (DiffTableHeader, error) {
	if !strings.HasSuffix(url, ".json") {
		return DiffTableHeader{}, fmt.Errorf("only .json format url is supported, sorry :(")
	}
	dth := DiffTableHeader{}
	common.FetchJson(url, &dth)
	return dth, nil
}

// Add a difficult table header(or say, meta data) and related song data to disk.
func (header *DiffTableHeader) SaveDiffTableHeader() error {
	// 1. Validation
	if arr, err := QueryDiffTableHeaderByName(header.Name); err != nil || len(arr) > 0 {
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
	db.Begin()
	if _, err := db.NamedExec(`INSERT INTO difftable_header(data_url, data_location, last_update, name, symbol, alias) VALUES (:data_url, :data_location, :last_update, :name, :symbol, :alias)`, header); err != nil {
		db.MustBegin().Rollback()
		return err
	}
	// 4. Try save course info into database
	if err := saveCourseInfoFromTableHeader(db, *header); err != nil {
		db.MustBegin().Rollback()
		return err
	}
	db.MustBegin().Commit()
	return nil
}

// Delete a difficult table header and its data file
func (header *DiffTableHeader) DeleteDiffTableHeader() error {
	db := common.OpenDB()
	defer db.Close()

	// 1) Try remove the data file, ignore any error
	os.Remove(header.DataLocation)
	// 2) Remove header from database
	_, err := db.Exec("DELETE FROM difftable_header WHERE id=?", header.Id)
	return err
}

// Return difficult table header's data json file name
func (header *DiffTableHeader) getDataJsonFileName() string {
	return header.Name + ".json"
}

// Fetch all data from sqlite
func QueryAllDiffTableHeader() ([]DiffTableHeader, error) {
	db := common.OpenDB()
	defer db.Close()
	var ret []DiffTableHeader
	err := db.Select(&ret, "SELECT * FROM difftable_header")
	return ret, err
}

// Query by name or alias
func QueryDiffTableHeaderByName(name string) ([]DiffTableHeader, error) {
	db := common.OpenDB()
	defer db.Close()
	var ret []DiffTableHeader
	err := db.Select(&ret, "SELECT * FROM difftable_header WHERE name=? OR alias=?", name, name)
	return ret, err
}

// Simple choose wrapper of QueryDifficultTableHeaderByName
func QueryDiffTableHeaderByNameWithChoices(name string) (DiffTableHeader, error) {
	dthArr, err := QueryDiffTableHeaderByName(name)
	if err != nil {
		return DiffTableHeader{}, err
	}
	return openDiffTableChooseTui(dthArr)
}

// Like QueryDiffTableHeaderByNameWithChoices, but without query
func AllDiffTableHeaderWithChoices() (DiffTableHeader, error) {
	dthArr, err := QueryAllDiffTableHeader()
	if err != nil {
		return DiffTableHeader{}, err
	}
	return openDiffTableChooseTui(dthArr)
}

// Simple choose wrapper for difftable_header
func openDiffTableChooseTui(dthArr []DiffTableHeader) (DiffTableHeader, error) {
	if len(dthArr) == 0 {
		return DiffTableHeader{}, fmt.Errorf("no table data")
	}
	choices := make([]string, 0)
	for _, v := range dthArr {
		choices = append(choices, v.String())
	}
	i := choose.OpenChooseTui(choices, "Multiple table matched with %s, please choose one:")
	return dthArr[i], nil
}

func (header *DiffTableHeader) SyncDifficultTable() error {
	// Sync data.json, propagate any error
	// This might lead to a corrupted file, but its ok since we can call sync command again to overwrite it
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
