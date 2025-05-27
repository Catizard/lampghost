package entity

import "gorm.io/gorm"

const (
	TASK_PREPARE = iota
	TASK_DOWNLOAD
	TASK_SUCCESS
	TASK_ERROR
)

type DownloadTask struct {
	gorm.Model
	URL string
	// NOTE: This field is designed to be a pointer due to gorm's update strategy
	Status               *int
	IntermediateFilePath string
	TaskName             *string
	DownloadSize         int64
	ContentLength        int64
	ErrorMessage         string
}
