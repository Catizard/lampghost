package entity

import "gorm.io/gorm"

type CourseInfo struct {
	gorm.Model
	SourceTableId string
	Name       string   
	Md5        []string 
	SourceId   int      
	SourceName string   
	Constraint []string 	
	// This field's only purpose is to store value in database
	// Since you cannot directly store array in database
	Md5s    string 
	Sha256s string // Can be seen as a mapping from md5s
}

func (CourseInfo) TableName() string {
	return "course_info"
}
