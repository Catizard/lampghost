package service

import (
	"fmt"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"

	. "github.com/samber/lo"
)

// Basic query function for diff_table_data table
//
// Difficult table's data might be lack of sha256 or md5. Therefore, we need to
// do repairment manually. This is why this function returns Dto instead of Entity.
func findDiffTableDataByID(tx *gorm.DB, ID uint) (*dto.DiffTableDataDto, error) {
	var data *entity.DiffTableData
	if err := tx.Find(&data, ID).Error; err != nil {
		return nil, eris.Wrap(err, "failed to find diff_table_data by id")
	}
	defaultCache, err := queryDefaultSongHashCache(tx)
	if err != nil {
		return nil, eris.Wrap(err, "failed to query default cache")
	}
	return dto.NewDiffTableDataDtoWithCache(data, defaultCache), nil
}

// Basic count query function for diff_table_data table
func selectDiffTableDataCount(tx *gorm.DB, filter *vo.DiffTableDataVo) (int64, error) {
	partial := tx.Model(&entity.DiffTableData{})
	if filter != nil {
		partial = partial.Where(filter.Entity())
	}
	var count int64
	if err := partial.Count(&count).Error; err != nil {
		return 0, eris.Wrap(err, "failed to query")
	}
	return count, nil
}

// Basic query function for diff_table_data table
func findDiffTableDataList(tx *gorm.DB, filter *vo.DiffTableDataVo) ([]*dto.DiffTableDataDto, int, error) {
	var contents []*entity.DiffTableData
	partial := tx.Model(&entity.DiffTableData{}).Scopes(scopeDiffTableDataFilter(filter))

	if err := partial.Debug().Find(&contents).Error; err != nil {
		return nil, 0, eris.Wrap(err, "failed to query")
	}

	if filter != nil && filter.Pagination != nil {
		count, err := selectDiffTableDataCount(tx, filter)
		if err != nil {
			return nil, 0, eris.Wrap(err, "failed to query")
		}
		filter.Pagination.PageCount = calcPageCount(count, filter.Pagination.PageSize)
	}

	defaultCache, err := queryDefaultSongHashCache(tx)
	if err != nil {
		return nil, 0, eris.Wrap(err, "failed to query default song hash cache")
	}
	ret := Map(contents, func(content *entity.DiffTableData, _ int) *dto.DiffTableDataDto {
		return dto.NewDiffTableDataDtoWithCache(content, defaultCache)
	})
	return ret, len(ret), nil
}

// Query difficult table tags, for building tag like "SL5", "ST11" for play logs
/*
select dh.name as table_name, dd."level" as table_level, dh.symbol as table_symbol, dh.tag_color as table_tag_color, dh.tag_text_color as table_tag_text_color
from difftable_data dd
left join difftable_header dh on dd.header_id = dh.id
where dd.md5 in ("176c2b2db4efd66cf186caae7923d477")
*/
func queryDiffTableTag(tx *gorm.DB, filter *vo.DiffTableDataVo) ([]*dto.DiffTableTagDto, int, error) {
	if filter == nil {
		return nil, 0, eris.New("queryDiffTableTag: filter cannot be nil")
	}
	fields := `difftable_data.md5, dh.name as table_name, difftable_data."level" as table_level, dh.symbol as table_symbol, dh.tag_color as table_tag_color, dh.tag_text_color as table_tag_text_color`
	var out []*dto.DiffTableTagDto
	partial := tx.Model(&entity.DiffTableData{}).
		Select(fields).
		Joins("left join difftable_header dh on difftable_data.header_id = dh.id").
		Where(filter.Entity()).
		Where("dh.no_tag_build = 0")
	if len(filter.Md5s) > 0 {
		partial = partial.Where("difftable_data.md5 in (?)", filter.Md5s)
	}
	if err := partial.Debug().Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, len(out), nil
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
	partial := tx.Table("difftable_data").Scopes(scopeDiffTableDataFilter(filter))

	partial = partial.Joins("left join (select id, sha256, md5, sub_title from rival_song_data group by md5) rsd on difftable_data.md5 = rsd.md5")
	partial = partial.Joins(`left join (
    select rsl.clear as Lamp, rsl.PlayCount, rsl.minbp as MinBP, rsl.sha256
    from (
      select rsl.clear, rsl.minbp, ROW_NUMBER() OVER w as rn, COUNT(1) OVER w as PlayCount, rsl.rival_id, rsl.sha256
      from rival_score_log rsl
			where rsl.rival_id = ?
      WINDOW w AS (PARTITION BY rsl.sha256 ORDER BY rsl.clear desc, rsl.minbp asc)
    ) as rsl
    where rsl.rn = 1
	) as rsl on rsl.sha256 = rsd.sha256`, filter.RivalID)
	partial = partial.Joins(`left join (
		select max(record_time) as record_time, sha256
		from rival_score_data_log
		where rival_id = ?
		group by sha256
	) as rsdl on rsdl.sha256 = rsl.sha256`, filter.RivalID)
	if filter.GhostRivalID > 0 {
		endRecordTime := maximumEndRecordTime
		if !filter.EndGhostRecordTime.IsZero() {
			endRecordTime = filter.EndGhostRecordTime
		}
		partial = partial.Joins(`left join (
			select rsl.clear as Lamp, rsl.PlayCount, rsl.minbp as MinBP, rsl.sha256
      from (
        select rsl.clear, rsl.minbp, ROW_NUMBER() OVER w as rn, COUNT(1) OVER w as PlayCount, rsl.rival_id, rsl.sha256
        from rival_score_log rsl
				where rsl.rival_id = ?
        WINDOW w AS (PARTITION BY rsl.sha256 ORDER BY rsl.clear desc, rsl.minbp asc)
      ) as rsl
      where rsl.rn = 1 and rsl.record_time <= ?
		) as ghost_rsl on ghost_rsl.sha256 = rsd.sha256`, filter.GhostRivalID, endRecordTime)
	}

	fields := `
		difftable_data.*,
		rsd.sha256,
		rsl.Lamp as Lamp, rsl.PlayCount as PlayCount,
		(rsd.id is null) as data_lost,
    rsd.sub_title as sub_title,
		strftime("%s", rsdl.record_time) as LastPlayedTimestamp
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

// Common query scope for vo.DiffTableDataVo
func scopeDiffTableDataFilter(filter *vo.DiffTableDataVo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}
		moved := db.Where(filter.Entity())
		// Extra filters here
		moved = moved.Scopes(
			scopeInIDs(filter.IDs),
			scopeInHeaderIDs(filter.HeaderIDs),
			pagination(filter.Pagination),
		)
		if filter.SortOrder != nil && *filter.SortOrder != "" {
			moved = moved.Order(fmt.Sprintf("%s %s", *filter.SortBy, filter.GetOrder()))
		}
		if len(filter.Levels) > 0 {
			moved = moved.Where("level in (?)", filter.Levels)
		}
		return moved
	}
}
