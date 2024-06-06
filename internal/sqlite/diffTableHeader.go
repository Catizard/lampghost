package sqlite

import (
	"fmt"
	"strings"

	"github.com/Catizard/lampghost/internal/data/difftable"
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
	return findList(tx, filter)
}

func (s *DiffTableHeaderService) FindDiffTableHeaderById(id string) (*difftable.DiffTableHeader, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return findById(tx, id)
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

func (s *DiffTableHeaderService) UpdateDiffTableHeader(id string, upd difftable.DiffTableHeaderUpdate) (*difftable.DiffTableHeader, error) {
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

func (s *DiffTableHeaderService) DeleteDifftableheader(id string) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteDifftableHeader(tx, id); err != nil {
		return err
	}
	return nil
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

func findList(tx *Tx, filter difftable.DiffTableHeaderFilter) (_ []*difftable.DiffTableHeader, n int, err error) {
	where := []string{"1 = 1"}
	if v := filter.Id; v != nil {
		where = append(where, "id = :id")
	}
	if v := filter.Name; v != nil {
		where = append(where, "name = :name")
	}

	rows, err := tx.NamedQuery(`
		SELECT
			id,
			data_url,
			data_location,
			last_update,
			name,
			original_url,
			symbol,
			alias
		FROM difftable_header
		WHERE `+strings.Join(where, " AND "),
		filter)
	if err != nil {
		return nil, n, err
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

	return ret, n, nil
}

func findById(tx *Tx, id string) (*difftable.DiffTableHeader, error) {
	arr, _, err := findList(tx, difftable.DiffTableHeaderFilter{Id: &id})
	if err != nil {
		return nil, err
	} else if len(arr) == 0 {
		return nil, fmt.Errorf("panic: no data")
	}
	return arr[0], nil
}

func updateDiffTableHeader(tx *Tx, id string, upd difftable.DiffTableHeaderUpdate) (*difftable.DiffTableHeader, error) {
	dth, err := findById(tx, id)
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

func deleteDifftableHeader(tx *Tx, id string) error {
	if _, err := findById(tx, id); err != nil {
		return err
	}

	_, err := tx.Exec("DELETE FROM difftable_header WHERE id=?", id)
	return err
}
