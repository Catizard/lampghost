package service

import (
	"github.com/Catizard/lampghost_wails/internal/entity"
	"gorm.io/gorm"
)

const DEFAULT_BATCH_SIZE = 100

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

// db.Where("sha256 in ?", sha256s)
// Requirements: len(sha256) > 0
func scopeInSha256s(sha256s []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(sha256s) > 0 {
			return db.Where("sha256 in ?", sha256s)
		}
		return db
	}
}

// pagination with page & pageSize
//
// Forcements:
// 1) if page <= 0, set it to 1
// 2) if pageSize >= 100, set it to 100
// 3) if pageSize <= 0, set if to 10
//
// Warning:
// 1) This function would modify `pagination` to `return` the page param back
// 2) If pagination is nil, do nothing
func pagination(pagination *entity.Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination == nil {
			return db
		}
		pagination.Page = normalizePage(pagination.Page)
		pagination.PageSize = normalizePageSize(pagination.PageSize)
		offset := (pagination.Page - 1) * pagination.PageSize
		return db.Offset(offset).Limit(pagination.PageSize)
	}
}

func normalizePage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func normalizePageSize(pageSize int) int {
	if pageSize >= 100 {
		return 100
	}
	if pageSize <= 0 {
		return 10
	}
	return pageSize
}

// helper method for calculating `PageCount` field
func calcPageCount(count int64, pageSize int) int {
	// HACK: prevent us from a un-normalized pageSize
	pageSize = normalizePage(pageSize)
	return int((count + int64(pageSize) - 1) / int64(pageSize))
}
