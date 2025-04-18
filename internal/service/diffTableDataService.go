package service

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"gorm.io/gorm"
)

func findDiffTableDataByID(tx *gorm.DB, ID uint) (*entity.DiffTableData, error) {
	var data *entity.DiffTableData
	if err := tx.Find(&data, ID).Error; err != nil {
		return nil, err
	}
	if err := fixDiffTableDataHashField(tx, data); err != nil {
		return nil, err
	}
	return data, nil
}

func findDiffTableDataList(tx *gorm.DB, filter *vo.DiffTableDataVo) ([]*dto.DiffTableDataDto, int, error) {
	if filter == nil {
		var contents []*entity.DiffTableData
		if err := tx.Find(&contents).Error; err != nil {
			return nil, 0, err
		}

		ret := make([]*dto.DiffTableDataDto, len(contents))
		defaultCache, err := queryDefaultSongHashCache(tx)
		if err != nil {
			return nil, 0, err
		}
		for i := range contents {
			ret[i] = dto.NewDiffTableDataDtoWithCache(contents[i], defaultCache)
		}
		return ret, len(ret), nil
	}

	var contents []*entity.DiffTableData
	partial := tx.Where(filter.Entity()).Scopes(
		scopeInIDs(filter.IDs),
		scopeInHeaderIDs(filter.HeaderIDs),
		pagination(filter.Pagination),
	)

	if filter.SortOrder != nil && *filter.SortOrder != "" {
		partial.Order(fmt.Sprintf("%s %s", *filter.SortBy, filter.GetOrder()))
	}

	if err := partial.Debug().Find(&contents).Error; err != nil {
		return nil, 0, err
	}

	if filter.Pagination != nil {
		count, err := selectDiffTableDataCount(tx, filter)
		if err != nil {
			return nil, 0, err
		}
		filter.Pagination.PageCount = calcPageCount(count, filter.Pagination.PageSize)
	}

	ret := make([]*dto.DiffTableDataDto, len(contents))
	defaultCache, err := queryDefaultSongHashCache(tx)
	if err != nil {
		return nil, 0, err
	}
	for i := range contents {
		ret[i] = dto.NewDiffTableDataDtoWithCache(contents[i], defaultCache)
	}
	return ret, len(ret), nil
}

func selectDiffTableDataCount(tx *gorm.DB, filter *vo.DiffTableDataVo) (int64, error) {
	if filter == nil {
		var count int64
		if err := tx.Model(&entity.DiffTableData{}).Count(&count).Error; err != nil {
			return 0, err
		}
		return count, nil
	}
	var count int64
	if err := tx.Model(&entity.DiffTableData{}).Where(filter.Entity()).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Fix the hash field on difficult table data
//
// NOTE: This function uses default user's song data to build the cache
func fixDiffTableDataHashField(tx *gorm.DB, rawContents ...*entity.DiffTableData) error {
	cache, err := queryDefaultSongHashCache(tx)
	if err != nil {
		return err
	}
	for _, rawContent := range rawContents {
		rawContent.RepairHash(cache)
	}
	return nil
}
