package difftable

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Catizard/lampghost/internal/data/filter"
	"github.com/guregu/null/v5"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type DiffTableHeader struct {
	Id           int         `db:"id"`
	DataUrl      string      `json:"data_url" db:"data_url"`
	DataLocation string      `db:"data_location"`
	LastUpdate   string      `json:"last_update" db:"last_update"`
	Name         string      `json:"name" db:"name"`
	OriginalUrl  null.String `json:"original_url" db:"original_url"`
	Symbol       string      `json:"symbol" db:"symbol"`
	Alias        string      `json:"alias" db:"alias"`
	// Not database related fields
	Course [][]CourseInfo `json:"course"`
	Data   []DiffTableData
	// Maps Data by level field?
	// DataLevelMap []map[string][]*DiffTableData
}

func (header *DiffTableHeader) String() string {
	if len(header.Alias) > 0 {
		return fmt.Sprintf("%s(%s) [symbol=%s, url=%s]", header.Name, header.Alias, header.Symbol, header.Alias)
	}
	return fmt.Sprintf("%s [symbol=%s, url=%s]", header.Name, header.Symbol, header.Alias)
}

func (header *DiffTableHeader) LoadData() error {
	f, err := os.Open(header.DataLocation)
	if err != nil {
		return err
	}
	defer f.Close()
	body, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	var local []DiffTableData
	err = json.Unmarshal(body, &local)
	if err != nil {
		return err
	}
	header.Data = local
	return nil
}

type DiffTableHeaderService interface {
	// ---------- basic methods ----------
	FindDiffTableHeaderList(filter filter.Filter) ([]*DiffTableHeader, int, error)
	FindDiffTableHeaderById(id int) (*DiffTableHeader, error)
	InsertDiffTableHeader(dth *DiffTableHeader) error
	UpdateDiffTableHeader(id int, upd DiffTableHeaderUpdate) (*DiffTableHeader, error)
	DeleteDiffTableHeader(id int) error

	// Fetch and save difficult table header info from remote url
	//
	// Support url forms:
	// 1) .json file
	// 2) .html file
	FetchAndSaveDiffTableHeader(url string, alias string) (*DiffTableHeader, error)

	// Simple wrapper of FindDiffTableHeaderList
	// After query, open tui app and wait user select one
	FindDiffTableHeaderListWithChoices(msg string, filter filter.Filter) (*DiffTableHeader, error)
}

type DiffTableHeaderFilter struct {
	// Filtering fields
	Id       null.Int    `db:"id"`
	Name     null.String `db:"name"`
	NameLike null.String `db:"nameLike"`
}

func (f DiffTableHeaderFilter) GenerateWhereClause() string {
	where := []string{"1 = 1"}
	if v := f.Id; v.Valid {
		where = append(where, "id = :id")
	}
	if v := f.Name; v.Valid && len(v.ValueOrZero()) > 0 {
		where = append(where, "name = :name")
	}
	if v := f.NameLike; v.Valid && len(v.ValueOrZero()) > 0 {
		where = append(where, "name like concat('%', :nameLike, '%') or alias like concat('%', :nameLike, '%')")
	}
	return strings.Join(where, " AND ")
}

type DiffTableHeaderUpdate struct {
	Name   *string
	Symbol *string
}

func (d *DiffTableHeader) MergeUpdate(upd DiffTableHeaderUpdate) {
	d.Name = *upd.Name
	d.Symbol = *upd.Symbol
}

// struct DiffTableData represents one content of a difficult table
type DiffTableData struct {
	Artist    string
	Comment   string
	Level     string
	Lr2BmsId  string `json:"lr2_bmdid"`
	Md5       string
	NameDiff  string
	Title     string
	Url       string `json:"url"`
	UrlDiff   string `json:"url_diff"`
	Sha256    string `json:"-"`
	// Lamp status
	Lamp      int32  `json:"-"`
	GhostLamp int32  `json:"-"`
}
