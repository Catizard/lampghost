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

// pagination with page & pageSize
// Forcements:
// 1) if page <= 0, set it to 1
// 2) if pageSize >= 100, set it to 100
// 3) if pageSize <= 0, set if to 10
func pagination(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize >= 100 {
			pageSize = 100
		}
		if pageSize <= 0 {
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
