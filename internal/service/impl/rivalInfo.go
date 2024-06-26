package impl

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/sqlite"
	"github.com/Catizard/lampghost/internal/tui/choose"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
)

var _ rival.RivalInfoService = (*RivalInfoService)(nil)

type RivalInfoService struct {
	db *sqlite.DB
}

func NewRivalInfoService(db *sqlite.DB) *RivalInfoService {
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

func (s *RivalInfoService) UpdateRivalInfo(id int, updater rival.RivalInfoUpdater) (*rival.RivalInfo, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	nv, err := updateRivalInfo(tx, id, updater)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	return nv, err
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
	return tx.Commit()
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

func findRivalInfoList(tx *sqlite.Tx, filter rival.RivalInfoFilter) (_ []*rival.RivalInfo, _ int, err error) {
	rows, err := tx.NamedQuery("SELECT * FROM rival_info WHERE "+filter.GenerateWhereClause(), filter)
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

func findRivalInfoById(tx *sqlite.Tx, id int) (*rival.RivalInfo, error) {
	arr, _, err := findRivalInfoList(tx, rival.RivalInfoFilter{Id: null.IntFrom(int64(id))})
	if err != nil {
		return nil, err
	} else if len(arr) == 0 {
		return nil, fmt.Errorf("panic: no data")
	}
	return arr[0], nil
}

func insertRivalInfo(tx *sqlite.Tx, rivalInfo *rival.RivalInfo) error {
	_, err := tx.NamedExec(`INSERT INTO rival_info (name, score_log_path, song_data_path, lr2_user_data_path) VALUES (:name,:score_log_path,:song_data_path,:lr2_user_data_path)`, rivalInfo)
	return err
}

func updateRivalInfo(tx *sqlite.Tx, id int, updater rival.RivalInfoUpdater) (*rival.RivalInfo, error) {
	r, err := findRivalInfoById(tx, id)
	if err != nil {
		return r, err
	}
	if updater.GenerateSetClause() == "" {
		log.Warn("Updater has no field, skip")
		return r, nil
	}
	updater.MergeUpdate(r)
	_, err = tx.NamedExec("UPDATE rival_info SET " + updater.GenerateSetClause() + " WHERE id=:id", updater)
	return r, err
}

func deleteRivalInfo(tx *sqlite.Tx, id int) error {
	if _, err := findRivalInfoById(tx, id); err != nil {
		return err
	}

	_, err := tx.Exec("DELETE FROM rival_info WHERE id=?", id)
	return err
}
