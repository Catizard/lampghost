package difftable

import (
	"strings"

	"github.com/Catizard/lampghost/internel/common"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Represents one course
type CourseInfo struct {
	Id     int      `db:"id"`
	Name   string   `json:"name" db:"name"`
	Md5    []string `json:"md5"`
	Source string   `db:"source"`
	// This field's only purpose is to store value in database
	// Since you cannot directly store array in database
	Md5s    string `db:"md5s"`
	Sha256s string // Can be seen as a mapping from md5s
}

// Initialze course_info header
func InitCourseInfoTable() error {
	db := common.OpenDB()
	defer db.Close()
	_, err := db.Exec(`
	DROP TABLE IF EXISTS 'course_info';
	CREATE TABLE course_info (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		source TEXT NOT NULL,
		md5s TEXT NOT NULL
	);`)
	return err
}

// Save course info from difficult table's fetch result
// If there is already a content has the same name, md5s, source, skip it
func saveCourseInfoFromTableHeader(db *sqlx.DB, dth DiffTableHeader) error {
	// If there is no course...
	if dth.Course == nil || len(dth.Course) == 0 || len(dth.Course[0]) == 0 {
		return nil
	}
	// There is no need to care about race
	courseArr, err := QueryAllCourseInfo()
	if err != nil {
		return err
	}

	// TODO: hack on rawData to delete no good/no speed mutated courses
	rawData := dth.Course
	for _, arr := range rawData {
		for _, v := range arr {
			v.prepareBeforeSave(dth)

			skipFlag := false
			for _, p := range courseArr {
				if v.Name == p.Name && v.Md5s == p.Md5s && v.Source == p.Source {
					skipFlag = true
					break
				}
			}
			if skipFlag {
				continue
			}
			// OK, it's unique
			if _, err := db.NamedExec("INSERT INTO course_info(name, md5s, source) VALUES(:name, :md5s, :source)", &v); err != nil {
				return err
			}
		}
	}
	return nil
}

// Prerequisite before save function on CourseInfo
func (c *CourseInfo) prepareBeforeSave(dth DiffTableHeader) {
	c.Md5s = strings.Join(c.Md5, ",")
	c.Source = dth.Name
}

// Preqrequiste after read function on CourseInfo
func (c *CourseInfo) prepareAfterRead() {
	// Split md5s field back
	c.Md5 = strings.Split(c.Md5s, ",")
}

// Fetch all data
func QueryAllCourseInfo() ([]CourseInfo, error) {
	db := common.OpenDB()
	defer db.Close()
	var ret []CourseInfo
	err := db.Select(&ret, "SELECT * FROM course_info")
	for i := range ret {
		ret[i].prepareAfterRead()
	}
	return ret, err
}
