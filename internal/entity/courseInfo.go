package entity

import (
	"fmt"

	"gorm.io/gorm"
)

type CourseInfo struct {
	gorm.Model
	HeaderID uint
	Name     string
	Sha256s  string
	// NOTE: Never use md5 at data processing, use sha256 instead
	Md5s        string
	Constraints string
}

func (CourseInfo) TableName() string {
	return "course_info"
}

// Before create requirements:
//
//	1.Sha256s & Md5s, at least one of them should not be empty
//	2. Name should not be empty
//	3. HeaderID should > 0
func (courseInfo *CourseInfo) BeforeCreate(tx *gorm.DB) error {
	if courseInfo.Sha256s == "" && courseInfo.Md5s == "" {
		return fmt.Errorf("courseInfo: BeforeCreate: sha256s & md5s are both empty")
	}
	if courseInfo.Name == "" {
		return fmt.Errorf("courseInfo: BeforeCreate: name is empty")
	}
	if courseInfo.HeaderID == 0 {
		return fmt.Errorf("courseInfo: BeforeCreate: headerId is 0")
	}
	return nil
}
