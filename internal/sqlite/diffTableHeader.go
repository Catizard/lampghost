package sqlite

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Catizard/lampghost/internal/common"
	"github.com/Catizard/lampghost/internal/config"
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/tui/choose"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
)

// Ensure service implements interface
var _ difftable.DiffTableHeaderService = (*DiffTableHeaderService)(nil)

// Represents a service component for managing difficult table header
type DiffTableHeaderService struct {
	db *DB
}

func NewDiffTableHeaderService(db *DB) *DiffTableHeaderService {
	return &DiffTableHeaderService{db: db}
}

func (s *DiffTableHeaderService) FindDiffTableHeaderList(filter difftable.DiffTableHeaderFilter) ([]*difftable.DiffTableHeader, int, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findDiffTableHeaderList(tx, filter)
}

func (s *DiffTableHeaderService) FindDiffTableHeaderById(id int) (*difftable.DiffTableHeader, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return findDiffTableHeaderById(tx, id)
}

func (s *DiffTableHeaderService) InsertDiffTableHeader(dth *difftable.DiffTableHeader) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := insertDiffTableHeader(tx, dth); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *DiffTableHeaderService) UpdateDiffTableHeader(id int, upd difftable.DiffTableHeaderUpdate) (*difftable.DiffTableHeader, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ret, err := updateDiffTableHeader(tx, id, upd)
	if err != nil {
		return ret, err
	} else if err := tx.Commit(); err != nil {
		return ret, err
	}
	return ret, nil
}

func (s *DiffTableHeaderService) DeleteDiffTableHeader(id int) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteDiffTableHeader(tx, id); err != nil {
		return err
	}
	return nil
}

func (s *DiffTableHeaderService) FetchAndSaveDiffTableHeader(url string, alias string) (*difftable.DiffTableHeader, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	// 1) Prepare header
	dth, err := fetchDiffTableFromURL(url)
	if err != nil {
		return nil, err
	}
	dth.Alias = alias
	// 2) Validate
	if ex, err := existsByName(tx, dth.Name); err != nil {
		return nil, err
	} else if ex {
		log.Fatalf(`There is already a table named (or its alias matches) %s
		Use table sync command to sync table.
		Use table del command to delete table`, dth.Name)
	}
	// 3) Download data file
	// Setup data file location
	dth.DataLocation = config.JoinWorkingDirectory(dth.Name + ".json")
	downloadTableData(dth)
	// 4) Insert into database
	if err := insertDiffTableHeader(tx, dth); err != nil {
		return nil, err
	}
	// 5) Save course info
	if err := saveCourseInfoFromTableHeader(tx, dth); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return dth, nil
}

func (s *DiffTableHeaderService) FindDiffTableHeaderListWithChoices(msg string, filter difftable.DiffTableHeaderFilter) (*difftable.DiffTableHeader, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if arr, n, err := findDiffTableHeaderList(tx, filter); err != nil {
		return nil, err
	} else if n == 0 {
		return nil, fmt.Errorf("no table data")
	} else {
		choices := make([]string, 0)
		for _, v := range arr {
			choices = append(choices, v.String())
		}
		i := choose.OpenChooseTuiSkippable(choices, msg)
		return arr[i], nil
	}
}

func insertDiffTableHeader(tx *Tx, dth *difftable.DiffTableHeader) error {
	_, err := tx.NamedExec(`
		INSERT INTO difftable_header(
			data_url,
			data_location,
			last_update,
			name,
			symbol,
			alias
		)
		VALUES (:data_url, :data_location, :last_update, :name, :symbol, :alias)`, dth)
	return err
}

func findDiffTableHeaderList(tx *Tx, filter difftable.DiffTableHeaderFilter) (_ []*difftable.DiffTableHeader, _ int, err error) {
	where := []string{"1 = 1"}
	if v := filter.Id; v.Valid {
		where = append(where, "id = :id")
	}
	if v := filter.Name; v.Valid {
		where = append(where, "name = :name")
	}
	if v := filter.NameLike; v.Valid {
		where = append(where, "name like concat('%', :nameLike, '%') or alias like concat('%', :nameLike, '%')")
	}

	rows, err := tx.NamedQuery(`
		SELECT *
		FROM difftable_header
		WHERE `+strings.Join(where, " AND "),
		filter)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	ret := make([]*difftable.DiffTableHeader, 0)
	for rows.Next() {
		dth := &difftable.DiffTableHeader{}
		if err := rows.StructScan(dth); err != nil {
			return nil, 0, err
		}

		ret = append(ret, dth)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return ret, len(ret), nil
}

func findDiffTableHeaderById(tx *Tx, id int) (*difftable.DiffTableHeader, error) {
	arr, _, err := findDiffTableHeaderList(tx, difftable.DiffTableHeaderFilter{Id: null.IntFrom(int64(id))})
	if err != nil {
		return nil, err
	} else if len(arr) == 0 {
		return nil, fmt.Errorf("panic: no data")
	}
	return arr[0], nil
}

func updateDiffTableHeader(tx *Tx, id int, upd difftable.DiffTableHeaderUpdate) (*difftable.DiffTableHeader, error) {
	dth, err := findDiffTableHeaderById(tx, id)
	if err != nil {
		return dth, err
	}
	dth.MergeUpdate(upd)

	if _, err := tx.NamedExec(`
		UPDATE difftable_header
		SET name = :name,
			symbol = :symbol				
		WHERE id = :id
	`, upd); err != nil {
		return dth, err
	}
	return dth, nil
}

func deleteDiffTableHeader(tx *Tx, id int) error {
	if _, err := findDiffTableHeaderById(tx, id); err != nil {
		return err
	}

	_, err := tx.Exec("DELETE FROM difftable_header WHERE id=?", id)
	return err
}

func fetchDiffTableFromURL(url string) (*difftable.DiffTableHeader, error) {
	if !strings.HasSuffix(url, ".json") {
		return nil, fmt.Errorf("only .json format url is supported, sorry :(")
	}
	dth := &difftable.DiffTableHeader{}
	common.FetchJson(url, &dth)
	return dth, nil
}

func existsByName(tx *Tx, name string) (bool, error) {
	filter := difftable.DiffTableHeaderFilter{
		Name: null.StringFrom(name),
	}
	if _, n, err := findDiffTableHeaderList(tx, filter); err != nil {
		return false, err
	} else if n > 0 {
		return true, nil
	}
	return false, nil
}

// Download difficult table's data.json file
// TODO: Use database's transaction to protect download phase?
func downloadTableData(dth *difftable.DiffTableHeader) error {
	if len(dth.DataUrl) == 0 {
		return fmt.Errorf("downloadTableData: no data url")
	}
	// 1) Create data file
	file, err := os.Create(dth.DataLocation)
	if err != nil {
		return err
	}
	defer file.Close()
	// 2) Download and save
	// TODO: if dataUrl is not start with http?
	resp, err := http.Get(dth.DataUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(file, resp.Body)
	return nil
}
