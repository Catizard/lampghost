package entity

import "gorm.io/gorm"

type CustomCourse struct {
	gorm.Model
	CustomTableID uint
	Name          string
	Sha256s       string
	Md5s          string
	Constraints   string
	OrderNumber   int
}

func (CustomCourse) TableName() string {
	return "custom_course"
}
