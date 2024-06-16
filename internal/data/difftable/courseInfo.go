package difftable

import (
	_ "github.com/mattn/go-sqlite3"
)

// Represents one course
type CourseInfo struct {
	Id         int      `db:"id"`
	Name       string   `json:"name" db:"name"`
	Md5        []string `json:"md5"`
	SourceId   int      `db:"source_id"`
	SourceName string   `db:"source_name"`
	Constraint []string `json:"constraint"`
	// This field's only purpose is to store value in database
	// Since you cannot directly store array in database
	Md5s    string `db:"md5s"`
	Sha256s string // Can be seen as a mapping from md5s
}

type CourseInfoService interface {
	// ---------- basic methods ----------
	FindCourseInfoList(filter CourseInfoFilter) ([]*CourseInfo, int, error)
	FindCourseInfoById(id int) (*CourseInfo, error)
	InsertCourseInfo(courseInfo *CourseInfo) error
	DeleteCourseInfo(id int) error
}

type CourseInfoFilter struct {
	Id   *int
	Name *string
}
