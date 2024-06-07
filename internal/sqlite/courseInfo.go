package sqlite

import (
	"fmt"
	"strings"

	"github.com/Catizard/lampghost/internal/data/difftable"
)

var _ difftable.CourseInfoService = (*CourseInfoService)(nil)

// Represents a service component for managing course info
type CourseInfoService struct {
	db *DB
}

func NewCourseInfoService(db *DB) *CourseInfoService {
	return &CourseInfoService{db: db}
}

func (s *CourseInfoService) FindCourseInfoList(filter difftable.CourseInfoFilter) ([]*difftable.CourseInfo, int, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findCourseInfoList(tx, filter)
}

func (s *CourseInfoService) FindCourseInfoById(id int) (*difftable.CourseInfo, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return findCourseInfoById(tx, id)
}

func (s *CourseInfoService) InsertCourseInfo(courseInfo *difftable.CourseInfo) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := insertCourseInfo(tx, courseInfo); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *CourseInfoService) DeleteCourseInfo(id int) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteCourseInfo(tx, id); err != nil {
		return err
	}
	return nil
}

func findCourseInfoList(tx *Tx, filter difftable.CourseInfoFilter) (_ []*difftable.CourseInfo, n int, err error) {
	where := []string{"1 = 1"}
	if v := filter.Id; v != nil {
		where = append(where, "id = :id")
	}
	if v := filter.Name; v != nil {
		where = append(where, "name = :name")
	}

	rows, err := tx.NamedQuery(`
		SELECT *
		FROM course_info
		WHERE `+strings.Join(where, " AND "),
		filter)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	ret := make([]*difftable.CourseInfo, 0)
	for rows.Next() {
		dth := &difftable.CourseInfo{}
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

func findCourseInfoById(tx *Tx, id int) (*difftable.CourseInfo, error) {
	arr, _, err := findCourseInfoList(tx, difftable.CourseInfoFilter{Id: &id})
	if err != nil {
		return nil, err
	} else if len(arr) == 0 {
		return nil, fmt.Errorf("panic: no data")
	}
	return arr[0], nil
}

func insertCourseInfo(tx *Tx, courseInfo *difftable.CourseInfo) error {
	_, err := tx.NamedExec("INSERT INTO course_info(name, md5s, source) VALUES(:name, :md5s, :source)", courseInfo)
	return err
}

func deleteCourseInfo(tx *Tx, id int) error {
	if _, err := findCourseInfoById(tx, id); err != nil {
		return err
	}

	_, err := tx.Exec("DELETE FROM course_info WHERE id=?", id)
	return err
}
