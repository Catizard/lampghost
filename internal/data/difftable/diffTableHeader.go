package difftable

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Catizard/lampghost/internal/common"
	"github.com/Catizard/lampghost/internal/tui/choose"
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

type DiffTableHeaderService interface {
	// ---------- basic methods ----------
	FindDiffTableHeaderList(filter DiffTableHeaderFilter) ([]*DiffTableHeader, int, error)
	FindDiffTableHeaderById(id string) (*DiffTableHeader, error)
	InsertDiffTableHeader(dth *DiffTableHeader) error
	UpdateDiffTableHeader(id string, upd DiffTableHeaderUpdate) (*DiffTableHeader, error)
	DeleteDifftableheader(id string) error

	// Fetch and save difficult table header info from remote url
	//
	// Support url forms:
	// 1) .json file
	FetchAndSaveDiffTableHeader(url string, alias string) (*DiffTableHeader, error)
}

type DiffTableHeaderFilter struct {
	// Filtering fields
	Id   *string
	Name *string
}

type DiffTableHeaderUpdate struct {
	Name   *string
	Symbol *string
}

func (d *DiffTableHeader) MergeUpdate(upd DiffTableHeaderUpdate) {
	d.Name = *upd.Name
	d.Symbol = *upd.Symbol
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
	i := choose.OpenChooseTuiSkippable(choices, "Choose one table to delete:")
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
