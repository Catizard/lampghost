package service

import "gorm.io/gorm"

// This file defines the scopes shared between services file

// db.Where("ID in (?)", IDs)
// Requirements: len(IDs) > 0
func scopeInIDs(IDs []uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(IDs) > 0 {
			return db.Where("ID in ?", IDs)
		}
		return db
	}
}

// db.Where("folder_id in ?", folderIDs)
// Requirements: len(folderIDs) > 0
func scopeInFolderIDs(folderIDs []uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(folderIDs) > 0 {
			return db.Where("folder_id in ?", folderIDs)
		}
		return db
	}
}

// db.Where("header_id in ?", headerIDs)
// Requirements: len(headerIDs) > 0
func scopeInHeaderIDs(headerIDs []uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(headerIDs) > 0 {
			return db.Where("header_id in ?", headerIDs)
		}
		return db
	}
}
