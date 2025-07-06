package entity

import "gorm.io/gorm"

type CustomCourse struct {
	gorm.Model
	CustomTableID uint
	Name          string
	Constraints   string
	OrderNumber   int
}

func (CustomCourse) TableName() string {
	return "custom_course"
}
