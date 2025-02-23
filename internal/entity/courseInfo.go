package entity

import "gorm.io/gorm"

type CourseInfo struct {
	gorm.Model
	HeaderID    uint
	Name        string
	Md5s        string
	Constraints string
}

func (CourseInfo) TableName() string {
	return "course_info"
}
