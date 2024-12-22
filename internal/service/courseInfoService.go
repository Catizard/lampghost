package service

import (
	"gorm.io/gorm"
)

type CourseInfoService struct {
	db *gorm.DB
}

func NewCourseInfoSerivce(db *gorm.DB) *CourseInfoService {
	return &CourseInfoService{
		db: db,
	}
}
