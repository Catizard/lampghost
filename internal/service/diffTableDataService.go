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

// Extension method for findDiffTableDataList, returning some
// player related fields at the same time.
//
// This function's implement strategy is kind of different from findDiffTableDataList,
// all memory related operation is replaced with one single sql to do. This is
// mainly because this function is required to make player's related fields sortable.
//
// Requirements:
//  1. To reduce complexity, filter should not be nil
//  2. filter.RivalID should > 0, otherwise this function is meaningless
//
// TODO: This function makes me feel like I should migrate to query build library instead.
// Constructing this complex sql in ORM is not really easy
//
// Warning:
//  1. When sorting ghost rival's related fields, filter.GhostRivalID should > 0 or
//     the sql statement would be broken, this behaviour is intended.
//  2. There's no detailed log attached to the result returning, only the highest
//     lamp wound be returned
func findDiffTableDataListWithRival(tx *gorm.DB, filter *vo.DiffTableDataVo) ([]*dto.DiffTableDataDto, int, error) {
	if filter == nil {
		return nil, 0, fmt.Errorf("findDiffTableDataListWithRival: filter cannot be nil")
	}
	if filter.RivalID == 0 {
		return nil, 0, fmt.Errorf("findDiffTableDataListWithRival: filter.rivalID should not be zero")
	}
	var contents []*dto.DiffTableDataDto
	partial := tx.Table("difftable_data").Where(filter.Entity()).Scopes(
		scopeInIDs(filter.IDs),
		scopeInHeaderIDs(filter.HeaderIDs),
		pagination(filter.Pagination),
	)

	if filter.SortOrder != nil && *filter.SortOrder != "" {
		partial.Order(fmt.Sprintf("%s %s", *filter.SortBy, filter.GetOrder()))
	}

	partial = partial.Joins("left join rival_song_data rsd on difftable_data.md5 = rsd.md5")
	partial = partial.Joins(`left join (
		select max(clear) as Lamp, count(1) as PlayCount, rsl.sha256
		from rival_score_log rsl
		where rsl.rival_id = ?
		group by rsl.sha256
	) as rsl on rsl.sha256 = rsd.sha256`, filter.RivalID)
	if filter.GhostRivalID > 0 {
		// TODO: How to do this???
		if !filter.EndGhostRecordTime.IsZero() {
			partial = partial.Joins(`left join (
			  select max(clear) as Lamp, count(1) as PlayCount, rsl.sha256
			  from rival_score_log rsl
			  where rsl.rival_id = ? and rsl.record_time <= ?
			  group by rsl.sha256
		    ) as ghost_rsl on ghost_rsl.sha256 = rsd.sha256`, filter.GhostRivalID, filter.EndGhostRecordTime)
		} else {
			partial = partial.Joins(`left join (
			  select max(clear) as Lamp, count(1) as PlayCount, rsl.sha256
			  from rival_score_log rsl
			  where rsl.rival_id = ?
			  group by rsl.sha256
		  ) as ghost_rsl on ghost_rsl.sha256 = rsd.sha256`, filter.GhostRivalID)
		}
	}

	fields := `
		difftable_data.*,
		rsl.Lamp as Lamp, rsl.PlayCount as PlayCount
	`
	if filter.GhostRivalID > 0 {
		fields = fields + ", ghost_rsl.Lamp as GhostLamp, ghost_rsl.PlayCount as GhostPlayCount"
	}

	if err := partial.Debug().Select(fields).Find(&contents).Error; err != nil {
		return nil, 0, err
	}

	if filter.Pagination != nil {
		count, err := selectDiffTableDataCount(tx, filter)
		if err != nil {
			return nil, 0, err
		}
		filter.Pagination.PageCount = calcPageCount(count, filter.Pagination.PageSize)
	}

	return contents, len(contents), nil
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
