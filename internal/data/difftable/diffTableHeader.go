package difftable

import (
	"fmt"

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
	FindDiffTableHeaderById(id int) (*DiffTableHeader, error)
	InsertDiffTableHeader(dth *DiffTableHeader) error
	UpdateDiffTableHeader(id int, upd DiffTableHeaderUpdate) (*DiffTableHeader, error)
	DeleteDiffTableHeader(id int) error

	// Fetch and save difficult table header info from remote url
	//
	// Support url forms:
	// 1) .json file
	FetchAndSaveDiffTableHeader(url string, alias string) (*DiffTableHeader, error)

	// Simple wrapper of FindDiffTableHeaderList
	// After query, open tui app and wait user select one
	FindDiffTableHeaderListWithChoices(msg string, filter DiffTableHeaderFilter) (*DiffTableHeader, error)
}

type DiffTableHeaderFilter struct {
	// Filtering fields
	Id   *int
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
