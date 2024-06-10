package sqlite

import (
	"fmt"
	"strings"

	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/tui/choose"
)

var _ rival.RivalInfoService = (*RivalInfoService)(nil)

type RivalInfoService struct {
	db *DB
}

func NewRivalInfoService(db *DB) *RivalInfoService {
	return &RivalInfoService{db: db}
}

func (s *RivalInfoService) FindRivalInfoList(filter rival.RivalInfoFilter) ([]*rival.RivalInfo, int, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findRivalInfoList(tx, filter)
}

func (s *RivalInfoService) FindRivalInfoById(id int) (*rival.RivalInfo, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return findRivalInfoById(tx, id)
}

func (s *RivalInfoService) InsertRivalInfo(dth *rival.RivalInfo) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := insertRivalInfo(tx, dth); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RivalInfoService) DeleteRivalInfo(id int) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteRivalInfo(tx, id); err != nil {
		return err
	}
	return nil
}

func (s *RivalInfoService) ChooseOneRival(msg string, filter rival.RivalInfoFilter) (*rival.RivalInfo, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if arr, n, err := findRivalInfoList(tx, filter); err != nil {
		return nil, err
	} else if n == 0 {
		return nil, fmt.Errorf("no rival data")
	} else {
		choices := make([]string, 0)
		for _, v := range arr {
			choices = append(choices, v.String())
		}
		i := choose.OpenChooseTuiSkippable(choices, msg)
		return arr[i], nil
	}
}

func findRivalInfoList(tx *Tx, filter rival.RivalInfoFilter) (_ []*rival.RivalInfo, n int, err error) {
	where := []string{"1 = 1"}
	if v := filter.Id; v != nil {
		where = append(where, "id = :id")
func findRivalInfoList(tx *Tx, filter rival.RivalInfoFilter) (_ []*rival.RivalInfo, _ int, err error) {
	}
	if v := filter.Name; v != nil {
		where = append(where, "name = :name")
	}

	rows, err := tx.NamedQuery(`
		SELECT *
		FROM rival_info
		WHERE `+strings.Join(where, " AND "),
		filter)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	ret := make([]*rival.RivalInfo, 0)
	for rows.Next() {
		r := &rival.RivalInfo{}
		if err := rows.StructScan(r); err != nil {
			return nil, 0, err
		}

		ret = append(ret, r)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return ret, len(ret), nil
}

func findRivalInfoById(tx *Tx, id int) (*rival.RivalInfo, error) {
	arr, _, err := findRivalInfoList(tx, rival.RivalInfoFilter{Id: &id})
	if err != nil {
		return nil, err
	} else if len(arr) == 0 {
		return nil, fmt.Errorf("panic: no data")
	}
	return arr[0], nil
}

func insertRivalInfo(tx *Tx, rivalInfo *rival.RivalInfo) error {
	_, err := tx.NamedExec(`INSERT INTO rival_info (name, score_log_path, song_data_path) VALUES (:name,:score_log_path,:song_data_path)`, rivalInfo);
	return err
}

func deleteRivalInfo(tx *Tx, id int) error {
	if _, err := findRivalInfoById(tx, id); err != nil {
		return err
	}

	_, err := tx.Exec("DELETE FROM rival_info WHERE id=?", id)
	return err
}