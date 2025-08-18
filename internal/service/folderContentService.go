package service

import (
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

func findFolderContentByID(tx *gorm.DB, id uint) (out *entity.FolderContent, err error) {
	err = tx.First(&out, id).Error
	return
}

func findFolderContentList(tx *gorm.DB, filter *vo.FolderContentVo) ([]*entity.FolderContent, int, error) {
	var contents []*entity.FolderContent
	if err := tx.Debug().Scopes(scopeFolderContentFilter(filter)).Find(&contents).Error; err != nil {
		return nil, 0, err
	}
	return contents, len(contents), nil
}

// Extension method for findFolderContentList, returning some
// player related fields and difficult table tags. This function is
// similar to findDiffTableDataListWithRival.
//
// Also the difficult table tags are built in memory not sql, since
// one content could be related to multiple difficult tables, which is
// obviously conflict with pagination
//
// Requirements:
//  1. To reduce complexity, filter should not be nil
//  2. filter.RivalID should > 0, otherwise this function is meaningless
func findFolderContentListWithRival(tx *gorm.DB, filter *vo.FolderContentVo) ([]*dto.FolderContentDto, int, error) {
	if filter == nil {
		return nil, 0, eris.New("fidnFolderContentListWithRvai: filter cannot be nil")
	}
	if filter.RivalID == 0 {
		return nil, 0, eris.New("findFolderContentListWithRival: filter.RivalID should not be 0")
	}
	var contents []*dto.FolderContentDto
	partial := tx.Model(&entity.FolderContent{}).Scopes(
		scopeFolderContentFilter(filter),
		pagination(filter.Pagination),
	)
	partial = partial.Joins(`left join (
		select max(clear) as Lamp, count(1) as PlayCount, rsl.sha256, rsl.record_time
		from rival_score_log rsl
		where rsl.rival_id = ?
		group by rsl.sha256
	) as rsl on rsl.sha256 = folder_content.sha256`, filter.RivalID)

	fields := `
		folder_content.*,
		rsl.Lamp as Lamp, rsl.PlayCount as PlayCount,
    strftime("%s", rsl.record_time) as BestRecordTimestamp
	`

	if err := partial.Select(fields).Find(&contents).Error; err != nil {
		return nil, 0, eris.Wrap(err, "query folder_content")
	}

	if filter.Pagination != nil {
		count, err := selectFolderContentCount(tx, filter)
		if err != nil {
			return nil, 0, eris.Wrap(err, "query folder_cotent")
		}
		filter.Pagination.PageCount = calcPageCount(count, filter.Pagination.PageSize)
	}
	return contents, len(contents), nil
}

func selectFolderContentCount(tx *gorm.DB, filter *vo.FolderContentVo) (int64, error) {
	moved := tx.Model(&entity.FolderContent{})
	var count int64
	if err := moved.Scopes(scopeFolderContentFilter(filter)).Count(&count).Error; err != nil {
		return 0, eris.Wrap(err, "query folder_content")
	}
	return count, nil
}

// Common query scope for vo.FolderContentVo
func scopeFolderContentFilter(filter *vo.FolderContentVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Add extra filter here
		moved = moved.Scopes(
			scopeInIDs(filter.IDs),
			scopeInFolderIDs(filter.FolderIDs),
		)
		return moved
	}
}
