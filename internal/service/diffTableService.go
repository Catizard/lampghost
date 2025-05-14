package service

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Catizard/bmstable"
	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/rotisserie/eris"
	. "github.com/samber/lo"
	"gorm.io/gorm"
)

type DiffTableService struct {
	db *gorm.DB
}

func NewDiffTableService(db *gorm.DB) *DiffTableService {
	return &DiffTableService{
		db: db,
	}
}

// Insert one difficult table by providing its url
//
// Requirements:
//
//	1.difficult table's url must be unique
//	2.difficult table's data_url must be unique
func (s *DiffTableService) AddDiffTableHeader(param *vo.DiffTableHeaderVo) (*entity.DiffTableHeader, error) {
	url := strings.TrimSpace(param.HeaderUrl)
	log.Debugf("[DiffTableService] calling AddDiffTableHeader with url: %s", url)
	if isDuplicated, err := queryDiffTableHeaderExistence(s.db, &entity.DiffTableHeader{HeaderUrl: url}); isDuplicated || err != nil {
		if err != nil {
			return nil, eris.Wrap(err, "failed to query duplicate header url")
		}
		return nil, eris.Errorf("add difficult table header failed: header_url[%s] is duplicated", url)
	}
	rawHeader, err := bmstable.ParseFromURL(url)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to parse difficult table url: %s", url)
	}
	if rawHeader.DataURL == "" {
		return nil, eris.Errorf("assert: data_url is empty")
	}
	if isDuplicated, err := queryDiffTableHeaderExistence(s.db, &entity.DiffTableHeader{DataUrl: rawHeader.DataURL}); isDuplicated || err != nil {
		if err != nil {
			return nil, eris.Wrap(err, "failed to query duplicate data url")
		}
		return nil, eris.Errorf("add difficult table header failed: data_url[%s] is duplicated", url)
	}

	headerEntity := entity.NewDiffTableHeaderFromImport(&rawHeader)
	// HACK: Inherit some fields from passing parameter
	headerEntity.TagColor = param.TagColor
	headerEntity.TagTextColor = param.TagTextColor

	// Transaction begins from here
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		// (1) difficult table header
		if err := tx.Create(headerEntity).Error; err != nil {
			return eris.Wrap(err, "failed to insert new header")
		}
		// (2) difficult related course contents
		courses := Map(rawHeader.Courses, func(importCourse bmstable.CourseInfo, _ int) *entity.CourseInfo {
			course := entity.NewCourseInfoFromImport(&importCourse)
			course.HeaderID = headerEntity.ID
			return course
		})
		if len(courses) > 0 {
			if err := delCourseInfo(tx, &vo.CourseInfoVo{HeaderID: headerEntity.ID}); err != nil {
				return eris.Wrap(err, "failed to delete previous courses")
			}
			if err := addBatchCourseInfo(tx, courses); err != nil {
				return eris.Wrap(err, "failed to insert new courses")
			}
		}
		// (3) difficult table concreate contents
		diffTableData := Map(rawHeader.Contents, func(importData bmstable.DifficultTableData, _ int) *entity.DiffTableData {
			data := entity.NewDiffTableDataFromImport(&importData)
			data.HeaderID = headerEntity.ID
			return data
		})
		if err := tx.Unscoped().Where("header_id = ?", headerEntity.ID).Delete(&entity.DiffTableData{}).Error; err != nil {
			return eris.Wrap(err, "failed to delete previous difficult table data")
		}
		if err := tx.CreateInBatches(&diffTableData, DEFAULT_BATCH_SIZE).Error; err != nil {
			return eris.Wrap(err, "failed to insert new difficult table data")
		}
		return nil
	}); err != nil {
		return nil, eris.Wrap(err, "transaction failed")
	}

	return headerEntity, nil
}

func (s *DiffTableService) UpdateDiffTableHeader(param *vo.DiffTableHeaderVo) error {
	if param == nil {
		return eris.Errorf("update: param cannot be nil")
	}
	if param.ID == 0 {
		return eris.Errorf("update: ID cannot be 0")
	}
	return s.db.Debug().Updates(param.Entity()).Error
}

// Query all difficult table datas
//
// Returns difficult header and its contents
func (s *DiffTableService) FindDiffTableHeaderList(filter *vo.DiffTableHeaderVo) ([]*dto.DiffTableHeaderDto, int, error) {
	headers, _, err := findDiffTableHeaderList(s.db, filter)
	if err != nil {
		return nil, 0, err
	}
	headerIds := Map(headers, func(item *entity.DiffTableHeader, index int) uint {
		return item.ID
	})
	contents, _, err := findDiffTableDataList(s.db, &vo.DiffTableDataVo{HeaderIDs: headerIds})
	if err != nil {
		return nil, 0, err
	}

	ret := Map(headers, func(header *entity.DiffTableHeader, _ int) *dto.DiffTableHeaderDto {
		return dto.NewDiffTableHeaderDto(header, Filter(contents, func(content *dto.DiffTableDataDto, _ int) bool {
			return content.HeaderID == header.ID
		}))
	})
	return ret, len(ret), nil
}

// Extend function for FindDiffTableHeaderList
//
// Adds player related field (e.g PlayCount, Lamp status)
func (s *DiffTableService) FindDiffTableHeaderListWithRival(rivalID uint) ([]*dto.DiffTableHeaderDto, int, error) {
	headers, _, err := s.FindDiffTableHeaderList(nil)
	if err != nil {
		return nil, 0, err
	}
	sha256ScoreLogsMap, err := findRivalScoreLogSha256Map(s.db, &vo.RivalScoreLogVo{RivalId: rivalID})
	if err != nil {
		return nil, 0, err
	}
	for _, header := range headers {
		mergeRivalRelatedData(sha256ScoreLogsMap, header.Contents, false)
	}
	return headers, len(headers), nil
}

// Query difficult table data as tree
//
// Example result:
//  1. Satelite:
//     1.1 satelite0
//     1.2 satelite1
//     1.3 satelite2
//     ... ....
//  2. BMS Insane table
//     ...
func (s *DiffTableService) FindDiffTableHeaderTree(filter *vo.DiffTableHeaderVo) ([]*dto.DiffTableHeaderDto, int, error) {
	// NOTE: Don't call s.FindDiffTableHeaderList, call findDiffTableHeaderList instead
	rawHeaders, _, err := findDiffTableHeaderList(s.db, filter)
	if err != nil {
		return nil, 0, err
	}

	if len(rawHeaders) == 0 {
		return make([]*dto.DiffTableHeaderDto, 0), 0, nil
	}

	headerIDs := make([]uint, 0)
	for _, header := range rawHeaders {
		headerIDs = append(headerIDs, header.ID)
	}

	pairs, err := queryRelatedLevelByIDS(s.db, headerIDs)
	if err != nil {
		return nil, 0, err
	}

	headers := make([]*dto.DiffTableHeaderDto, 0)
	for _, header := range rawHeaders {
		headerDto := dto.NewDiffTableHeaderDto(header, nil)
		for _, pair := range pairs {
			if header.ID == pair.header_id {
				headerDto.Children = append(headerDto.Children, *dto.NewLevelChildNode(
					header.ID,
					fmt.Sprintf("%s%s", header.Symbol, pair.level),
					pair.level,
				))
			}
		}
		headers = append(headers, headerDto)
	}

	for headerInx, header := range headers {
		headers[headerInx].Children = sortHeadersByLevel(header.Children, header.UnjoinedLevelOrder)
	}

	return headers, len(headers), nil
}

// Extend function for FindDiffTableHeaderTree
//
// Adds player related field (e.g PlayCount, Lamp status)
func (s *DiffTableService) FindDiffTableHeaderTreeWithRival(filter *vo.DiffTableHeaderVo) ([]*dto.DiffTableHeaderDto, int, error) {
	if filter == nil || filter.RivalID == 0 {
		return nil, 0, fmt.Errorf("FindDiffTableHeaderTreeWithRival: rival id is empty or 0")
	}
	// NOTE: Unlike FindDiffTableHeaderTree, this function must query complete table data
	// Therefore, there is no re-usable code
	rawHeaders, n, err := findDiffTableHeaderList(s.db, filter)
	if err != nil {
		return nil, 0, err
	}
	if n == 0 {
		return make([]*dto.DiffTableHeaderDto, 0), 0, nil
	}

	headerIDs := make([]uint, 0)
	for _, header := range rawHeaders {
		headerIDs = append(headerIDs, header.ID)
	}

	// Here, we must query the complete difficult table data
	// While FindDiffTableHeaderTree could only query levels
	rawDiffTableDataList, _, err := findDiffTableDataList(s.db, &vo.DiffTableDataVo{
		HeaderIDs: headerIDs,
	})
	if err != nil {
		return nil, 0, err
	}
	// headerID#level => [sha256]
	difftableSha256List := make(map[string][]string)
	headerIDMapsToLevelList := make(map[uint][]string)
	// headerID#level set
	dupLevelSet := make(map[string]any)
	for _, diffTableData := range rawDiffTableDataList {
		headerID := diffTableData.HeaderID
		key := fmt.Sprintf("%d#%s", headerID, diffTableData.Level)
		if _, ok := difftableSha256List[key]; !ok {
			difftableSha256List[key] = make([]string, 0)
		}
		difftableSha256List[key] = append(difftableSha256List[key], diffTableData.Sha256)
		if _, ok := headerIDMapsToLevelList[headerID]; !ok {
			headerIDMapsToLevelList[headerID] = make([]string, 0)
		}
		if _, ok := dupLevelSet[key]; !ok {
			headerIDMapsToLevelList[headerID] = append(headerIDMapsToLevelList[headerID], diffTableData.Level)
			dupLevelSet[key] = new(any)
		}
	}

	// Query player's maximum clear and group them by sha256
	scoreLogSha256Map, err := findRivalScoreLogSha256Map(s.db, &vo.RivalScoreLogVo{RivalId: filter.RivalID})
	if err != nil {
		return nil, 0, err
	}

	headers := make([]*dto.DiffTableHeaderDto, 0)
	for _, header := range rawHeaders {
		headerDto := dto.NewDiffTableHeaderDto(header, nil)
		for _, level := range headerIDMapsToLevelList[header.ID] {
			levelNode := dto.NewLevelChildNode(
				header.ID,
				fmt.Sprintf("%s%s", header.Symbol, level),
				level,
			)
			relatedSha256List := difftableSha256List[fmt.Sprintf("%d#%s", header.ID, level)]
			levelNode.SongCount = len(relatedSha256List)
			levelNode.LampCount = make(map[int]int)
			for _, sha256 := range relatedSha256List {
				if maximumLog, ok := scoreLogSha256Map[sha256]; ok {
					levelNode.LampCount[int(maximumLog[0].Clear)] += 1
				}
			}
			headerDto.Children = append(headerDto.Children, *levelNode)
		}
		headers = append(headers, headerDto)
	}

	for headerInx, header := range headers {
		headers[headerInx].Children = sortHeadersByLevel(header.Children, header.UnjoinedLevelOrder)
	}

	return headers, len(headers), nil
}

func (s *DiffTableService) DelDiffTableHeader(ID uint) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		var candidate entity.DiffTableHeader
		if err := tx.First(&candidate, ID).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("header_id = ?", candidate.ID).Delete(&entity.DiffTableData{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("header_id = ?", candidate.ID).Delete(&entity.CourseInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&entity.DiffTableHeader{}, candidate.ID).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Query specific difficult table's info
//
// Returns one header with detailed contents, ensure every content's hash field is compitable
func (s *DiffTableService) QueryDiffTableInfoByID(ID uint) (*dto.DiffTableHeaderDto, error) {
	return queryDiffTableInfoByID(s.db, ID)
}

// Extend function for QueryDiffTableInfoByID
//
// Adds player related field (e.g PlayCount, Lamp status)
func (s *DiffTableService) QueryDiffTableInfoByIDWithRival(ID uint, rivalID uint) (*dto.DiffTableHeaderDto, error) {
	header, err := s.QueryDiffTableInfoByID(ID)
	if err != nil {
		return nil, err
	}
	sha256ScoreLogsMap, err := findRivalScoreLogSha256Map(s.db, &vo.RivalScoreLogVo{RivalId: rivalID})
	if err != nil {
		return nil, err
	}
	mergeRivalRelatedData(sha256ScoreLogsMap, header.Contents, false)
	return header, nil
}

func (s *DiffTableService) QueryLevelLayeredDiffTableInfoByID(ID uint) (*dto.DiffTableHeaderDto, error) {
	return queryLevelLayeredDiffTableInfoById(s.db, ID)
}

// Query specific difficult table's one level data contents with player related field (e.g PlayCount, Lamp status...)
//
// Requirements:
//
//  1. Level & ID & RivalID should not be empty
//  2. The rival's data must be queryed with one sql, because these fields are sortable
func (s *DiffTableService) QueryDiffTableDataWithRival(filter *vo.DiffTableHeaderVo) ([]*dto.DiffTableDataDto, int, error) {
	if filter.Level == "" {
		return nil, 0, eris.Errorf("Level should not be empty")
	}
	if filter.ID <= 0 {
		return nil, 0, eris.Errorf("ID should > 0")
	}
	if filter.RivalID <= 0 {
		return nil, 0, eris.Errorf("RivalID should > 0")
	}

	var endGhostRecordTime time.Time
	if filter.GhostRivalTagID > 0 {
		tag, err := findRivalTagByID(s.db, filter.GhostRivalTagID)
		if err != nil {
			return nil, 0, eris.Wrap(err, "failed to query rival tag by id")
		}
		endGhostRecordTime = tag.RecordTime
	}
	return findDiffTableDataListWithRival(s.db, &vo.DiffTableDataVo{
		HeaderID:           filter.ID,
		Level:              filter.Level,
		Pagination:         filter.Pagination,
		SortBy:             filter.SortBy,
		SortOrder:          filter.SortOrder,
		RivalID:            filter.RivalID,
		GhostRivalID:       filter.GhostRivalID,
		EndGhostRecordTime: endGhostRecordTime,
	})
}

func (s *DiffTableService) BindDiffTableDataToFolder(diffTableDataID uint, folderIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		log.Debugf("before find difficult table data")
		songData, err := findDiffTableDataByID(tx, diffTableDataID)
		if err != nil {
			return err
		}
		log.Debugf("data: %v", songData)

		content := entity.FolderContent{
			Sha256: songData.Sha256,
			Md5:    songData.Md5,
			Title:  songData.Title,
		}

		log.Debugf("new folder content: %v", content)

		if err := bindSongToFolder(tx, content, folderIDs); err != nil {
			return err
		}
		return nil
	})
}

func (s *DiffTableService) UpdateHeaderOrder(headerIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return updateHeaderOrder(tx, headerIDs)
	})
}

func (s *DiffTableService) UpdateHeaderLevelOrders(updateParam *vo.DiffTableHeaderVo) error {
	if updateParam == nil {
		return fmt.Errorf("assert: UpdateHeaderLevelOrders: updateParam cannot be nil")
	}
	if updateParam.ID == 0 {
		return fmt.Errorf("assert: UpdateHeaderLevelOrders: updateParam.ID cannot be 0")
	}
	if updateParam.LevelOrders == "" {
		return fmt.Errorf("assert: UpdateHeaderLevelOrders: updateParam.LevelOrders cannot be empty")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		return updateHeaderLevelOrders(tx, updateParam.ID, updateParam.LevelOrders)
	})
}

// Query if there exists a header that satisfies the condition
func queryDiffTableHeaderExistence(tx *gorm.DB, filter *entity.DiffTableHeader) (bool, error) {
	var dupCount int64
	if err := tx.Model(&entity.DiffTableHeader{}).Where(filter).Count(&dupCount).Error; err != nil {
		return false, eris.Wrap(err, "failed to query difficult table header existence")
	}
	return dupCount > 0, nil
}

func queryLevelLayeredDiffTableInfoById(tx *gorm.DB, ID uint) (*dto.DiffTableHeaderDto, error) {
	header, err := queryDiffTableInfoByID(tx, ID)
	if err != nil {
		return nil, err
	}
	levels := make(map[string]any)
	levelLayeredContent := make(map[string][]*dto.DiffTableDataDto)
	for _, v := range header.Contents {
		if _, ok := levelLayeredContent[v.Level]; !ok {
			levelLayeredContent[v.Level] = make([]*dto.DiffTableDataDto, 0)
		}
		if _, ok := levels[v.Level]; !ok {
			levels[v.Level] = new(any)
		}
		levelLayeredContent[v.Level] = append(levelLayeredContent[v.Level], v)
	}

	sortedLevels := make([]string, 0)
	for level := range levels {
		sortedLevels = append(sortedLevels, level)
	}
	sortedLevels = sortLevels(sortedLevels, header.UnjoinedLevelOrder)
	return dto.NewLevelLayeredDiffTableHeaderDto(header.Entity(), sortedLevels, levelLayeredContent), nil
}

// Query one difficult table header and its related contents by header's ID
//
// Related contents would be attached on `Contents` field
func queryDiffTableInfoByID(tx *gorm.DB, ID uint) (*dto.DiffTableHeaderDto, error) {
	var header entity.DiffTableHeader
	if err := tx.First(&header, ID).Error; err != nil {
		return nil, err
	}
	contents, _, err := findDiffTableDataList(tx, &vo.DiffTableDataVo{HeaderID: ID})
	if err != nil {
		return nil, err
	}
	return dto.NewDiffTableHeaderDto(&header, contents), nil
}

// Merge player related data onto DiffTableDataDto (e.g PlayCount LampStatus...)
// TODO: We can actaully combine "query rival's related data" and "merge rival's data with DiffTableDataDto" two steps together
// The obstacle is mainly FindDiffTableHeaderListWithRival function, which requires redesign the data loading sequence
//
// This function would modify data in place rather than return a new array
func mergeRivalRelatedData(sha256ScoreLogsMap map[string][]*dto.RivalScoreLogDto, contents []*dto.DiffTableDataDto, isGhostRival bool) error {
	for i, content := range contents {
		if logs, ok := sha256ScoreLogsMap[content.Sha256]; ok {
			contents[i].PlayCount = len(logs)
			for _, log := range logs {
				if isGhostRival {
					contents[i].GhostLamp = max(content.GhostLamp, int(log.Clear))
				} else {
					contents[i].Lamp = max(content.Lamp, int(log.Clear))
				}
			}
		} else {
			contents[i].PlayCount = 0
		}
	}
	return nil
}

func findDiffTableHeaderList(tx *gorm.DB, filter *vo.DiffTableHeaderVo) ([]*entity.DiffTableHeader, int, error) {
	if filter == nil {
		var headers []*entity.DiffTableHeader
		if err := tx.Order("order_number").Find(&headers).Error; err != nil {
			log.Error("[DiffTableService] Find difftable header failed with %v", err)
			return nil, 0, err
		}
		return headers, len(headers), nil
	}

	var headers []*entity.DiffTableHeader
	if err := tx.Where(filter.Entity()).Order("order_number").Find(&headers).Error; err != nil {
		return nil, 0, err
	}
	return headers, len(headers), nil
}

func updateHeaderOrder(tx *gorm.DB, headerIDs []uint) error {
	if len(headerIDs) == 0 {
		return nil
	}
	for i, headerID := range headerIDs {
		entity := &entity.DiffTableHeader{}
		entity.ID = headerID
		if err := tx.Model(entity).Update("order_number", i).Error; err != nil {
			return err
		}
	}
	return nil
}

func updateHeaderLevelOrders(tx *gorm.DB, headerID uint, levelOrders string) error {
	entity := &entity.DiffTableHeader{}
	entity.ID = headerID
	return tx.Model(entity).Update("level_orders", levelOrders).Error
}

// Query multiple difficult table's related level list by header ids
// When only related level list are required, this function is cheapier than load whole data content
//
// NOTE: parameter IDs must not be empty or the sql structure isn't correct
// Returns a list of pair(header_id, level)
func queryRelatedLevelByIDS(tx *gorm.DB, IDs []uint) (ret []struct {
	header_id uint
	level     string
}, err error,
) {
	rows, err := tx.Raw(`select dd.header_id, dd."level"
		from difftable_data dd
		group by dd.header_id, dd."level"
		having dd.header_id in ? `, IDs).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var header_id uint
		var level string
		rows.Scan(&header_id, &level)
		ret = append(ret, struct {
			header_id uint
			level     string
		}{
			header_id: header_id,
			level:     level,
		})
	}

	return
}

// Sort difficult table headers by following rule:
//  1. if lhs and rhs are both defined in `preSortLevels`, return index(lhs) < index(rhs)
//  2. if one of them is not present,
//     2.1 if lhs and rhs are both number, return number(lhs) < number(rhs)
//     2.2 if one of them is not number, return lhs < rhs
func sortHeadersByLevel(headers []dto.DiffTableHeaderDto, preSortLevels []string) []dto.DiffTableHeaderDto {
	sorted := make([]dto.DiffTableHeaderDto, len(headers))
	copy(sorted, headers)
	slices.SortFunc(sorted, func(lhs dto.DiffTableHeaderDto, rhs dto.DiffTableHeaderDto) int {
		return levelComparator(lhs.Level, rhs.Level, preSortLevels)
	})
	return sorted
}

func sortLevels(levels []string, preSortLevels []string) []string {
	sorted := make([]string, len(levels))
	copy(sorted, levels)
	slices.SortFunc(sorted, func(lhs string, rhs string) int {
		return levelComparator(lhs, rhs, preSortLevels)
	})
	return sorted
}

// Compares two level string by following rule:
//  1. if lhs and rhs are both defined in `preSortLevels`, return index(lhs) < index(rhs)
//  2. if one of them is not present,
//     2.1 if lhs and rhs are both number, return number(lhs) < number(rhs)
//     2.2 if one of them is not number, return lhs < rhs
func levelComparator(lhs string, rhs string, preSortLevels []string) int {
	inxL := -1
	inxR := -1
	if len(preSortLevels) > 0 {
		inxL = slices.Index(preSortLevels, lhs)
		inxR = slices.Index(preSortLevels, rhs)
	}
	if inxL == -1 || inxR == -1 {
		ill, errL := strconv.Atoi(lhs)
		irr, errR := strconv.Atoi(rhs)
		if errL == nil && errR == nil {
			return ill - irr
		}
		if lhs < rhs {
			return -1
		} else if lhs > rhs {
			return 1
		}
		return 0
	}
	return inxL - inxR
}
