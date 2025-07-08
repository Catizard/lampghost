package entity

import "gorm.io/gorm"

type CustomCourseData struct {
	gorm.Model
	Sha256         string
	Md5            string
	CustomCourseID uint
	OrderNumber    int
}

func (CustomCourseData) TableName() string {
	return "custom_course_data"
}
