package sqlite

import (
	"fmt"
	"strings"

	"github.com/Catizard/lampghost/internal/data/difftable"
)

var (
	// TODO: Make shouldIgnore as a configurable option
	shouldIgnoreSpecialConstaints = true
	ignoreConstraints             = map[string]struct{}{
		"no_good":  {},
		"no_speed": {},
	}
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
		c := &difftable.CourseInfo{}
		if err := rows.StructScan(c); err != nil {
			return nil, 0, err
		}

		prepareAfterRead(c)
		ret = append(ret, c)
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
	prepareAfterRead(arr[0])
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

// Save course info from difficult table's fetch result
func saveCourseInfoFromTableHeader(tx *Tx, dth *difftable.DiffTableHeader) error {
	// If there is no course...
	if dth.Course == nil || len(dth.Course) == 0 || len(dth.Course[0]) == 0 {
		return nil
	}
	// There is no need to care about race
	courseArr, _, err := findCourseInfoList(tx, difftable.CourseInfoFilter{})
	if err != nil {
		return err
	}

	rawData := dth.Course
	for _, arr := range rawData {
		for _, v := range arr {
			prepareBeforeSave(&v, dth)

			skipFlag := false
			// Skip 1: There is a same course exists
			for _, p := range courseArr {
				if v.Name == p.Name && v.Md5s == p.Md5s && v.Source == p.Source {
					skipFlag = true
					break
				}
			}
			// Skip 2: Open ignore special constraints flag and it matches at least one
			if shouldIgnoreSpecialConstaints {
				for _, constraint := range v.Constraint {
					if _, ok := ignoreConstraints[constraint]; ok {
						skipFlag = true
					}
				}
			}
			if skipFlag {
				continue
			}
			// OK, it's unique
			if err := insertCourseInfo(tx, &v); err != nil {
				return err
			}
		}
	}
	return nil
}

// Prerequisite before save function on CourseInfo
func prepareBeforeSave(c *difftable.CourseInfo, dth *difftable.DiffTableHeader) {
	c.Md5s = strings.Join(c.Md5, ",")
	c.Source = dth.Name
}

// Preqrequiste after read function on CourseInfo
func prepareAfterRead(c *difftable.CourseInfo) {
	// Split md5s field back
	c.Md5 = strings.Split(c.Md5s, ",")
}
