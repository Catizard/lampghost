package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type DiffTableService struct {
	db                   *gorm.DB
	rivalSongDataService *RivalSongDataService
}

func NewDiffTableService(db *gorm.DB, rivalSongDataService *RivalSongDataService) *DiffTableService {
	return &DiffTableService{
		db:                   db,
		rivalSongDataService: rivalSongDataService,
	}
}

func (s *DiffTableService) AddDiffTableHeader(url string) (*entity.DiffTableHeader, error) {
	url = strings.TrimSpace(url)
	log.Debugf("[DiffTableService] calling AddDiffTableHeader with url: %s", url)
	if isDuplicated, err := queryDiffTableHeaderExistence(s.db, &entity.DiffTableHeader{HeaderUrl: url}); isDuplicated || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("add difficult table header failed: header_url[%s] is duplicated", url)
	}
	headerVo, err := fetchDiffTableFromURL(url)
	if err != nil {
		return nil, err
	}
	headerVo.HeaderUrl = url
	if headerVo.DataUrl == "" {
		return nil, fmt.Errorf("assert: header.DataUrl cannot be empty")
	}
	if !strings.HasSuffix(headerVo.DataUrl, ".json") {
		return nil, fmt.Errorf("assert: header.DataUrl must endes with .json")
	}
	log.Debugf("[DiffTableService] Got header data: %v", headerVo)
	if isDuplicated, err := queryDiffTableHeaderExistence(s.db, &entity.DiffTableHeader{DataUrl: headerVo.DataUrl}); isDuplicated || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("add difficult table header failed: data_url[%s] is duplicated", url)
	}

	// Transaction begins from here
	headerEntity := headerVo.Entity()
	var data []entity.DiffTableData
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		// (1) difficult table header
		if err := tx.Create(headerEntity).Error; err != nil {
			return err
		}
		// (2) difficult related course contents
		// NOTE: we have to do a custom decode here, some tables provide courses as a two-dimensional array
		// while some others provide a one-dimensional array instead
		if err := headerVo.ParseRawCourses(); err != nil {
			return err
		}
		if len(headerVo.Courses) > 0 {
			var courseData []entity.CourseInfo
			for _, courseInfoVo := range headerVo.Courses {
				courseInfo := courseInfoVo.Entity()
				courseInfo.HeaderID = headerEntity.ID
				courseData = append(courseData, *courseInfo)
			}
			if err := tx.Unscoped().Where("header_id = ?", headerEntity.ID).Delete(&entity.CourseInfo{}).Error; err != nil {
				return err
			}
			if err := tx.Create(&courseData).Error; err != nil {
				return err
			}
		}
		// (3) difficult table concreate contents
		if err := fetchJson(headerVo.DataUrl, &data); err != nil {
			return err
		}
		for i := range data {
			data[i].HeaderID = headerEntity.ID
		}
		if err := tx.Unscoped().Where("header_id = ?", headerEntity.ID).Delete(&entity.DiffTableData{}).Error; err != nil {
			return err
		}
		if err := tx.CreateInBatches(&data, DEFAULT_BATCH_SIZE).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	log.Infof("[DiffTableService] Inserted one header with %d contents", len(data))
	return headerEntity, nil
}

// Query all difficult table datas
//
// Returns difficult header and its contents
func (s *DiffTableService) FindDiffTableHeaderList(filter *vo.DiffTableHeaderVo) ([]dto.DiffTableHeaderDto, int, error) {
	headers, _, err := findDiffTableHeaderList(s.db, filter)
	if err != nil {
		return nil, 0, err
	}
	headerIds := make([]uint, len(headers))
	for i, header := range headers {
		headerIds[i] = header.ID
	}
	rawContents, _, err := findDiffTableDataList(s.db, &vo.DiffTableDataVo{HeaderIDs: headerIds})
	if err != nil {
		return nil, 0, err
	}

	ret := make([]dto.DiffTableHeaderDto, len(headers))
	for i, header := range headers {
		contents := make([]*dto.DiffTableDataDto, 0)
		for _, content := range rawContents {
			if content.HeaderID == header.ID {
				contents = append(contents, dto.NewDiffTableDataDto(content))
			}
		}
		ret[i] = *dto.NewDiffTableHeaderDto(header, contents)
	}

	return ret, len(ret), nil
}

// Extend function for FindDiffTableHeaderList
//
// Adds player related field (e.g PlayCount, Lamp status)
func (s *DiffTableService) FindDiffTableHeaderListWithRival(rivalID uint) ([]dto.DiffTableHeaderDto, int, error) {
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
func (s *DiffTableService) FindDiffTableHeaderTree(filter *vo.DiffTableHeaderVo) ([]dto.DiffTableHeaderDto, int, error) {
	// NOTE: Don't call s.FindDiffTableHeaderList, call findDiffTableHeaderList instead
	rawHeaders, _, err := findDiffTableHeaderList(s.db, filter)
	if err != nil {
		return nil, 0, err
	}

	if len(rawHeaders) == 0 {
		return make([]dto.DiffTableHeaderDto, 0), 0, nil
	}

	headerIDs := make([]uint, 0)
	for _, header := range rawHeaders {
		headerIDs = append(headerIDs, header.ID)
	}

	pairs, err := queryRelatedLevelByIDS(s.db, headerIDs)
	if err != nil {
		return nil, 0, err
	}

	headers := make([]dto.DiffTableHeaderDto, 0)
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
		headers = append(headers, *headerDto)
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
	header, err := queryDiffTableInfoByID(s.db, ID)
	if err != nil {
		return nil, err
	}
	rawContents, _, err := findDiffTableDataList(s.db, &vo.DiffTableDataVo{HeaderID: ID})
	if err != nil {
		return nil, err
	}
	contents := dto.NewDiffTableDataDtoArray(rawContents)
	return dto.NewDiffTableHeaderDto(header, contents), nil
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

func (s *DiffTableService) QueryLevelLayeredDiffTableInfoById(ID uint) (*dto.DiffTableHeaderDto, error) {
	header, err := s.QueryDiffTableInfoByID(ID)
	if err != nil {
		return nil, err
	}
	levels := make(map[string]interface{})
	levelLayeredContent := make(map[string][]*dto.DiffTableDataDto)
	for _, v := range header.Contents {
		if _, ok := levelLayeredContent[v.Level]; !ok {
			levelLayeredContent[v.Level] = make([]*dto.DiffTableDataDto, 0)
		}
		if _, ok := levels[v.Level]; !ok {
			levels[v.Level] = new(interface{})
		}
		levelLayeredContent[v.Level] = append(levelLayeredContent[v.Level], v)
	}

	sortedLevels := make([]string, 0)
	for level := range levels {
		sortedLevels = append(sortedLevels, level)
	}
	sort.Slice(sortedLevels, func(i, j int) bool {
		ll := sortedLevels[i]
		rr := sortedLevels[j]
		ill, errL := strconv.Atoi(ll)
		irr, errR := strconv.Atoi(rr)
		if errL == nil && errR == nil {
			return ill < irr
		}
		return ll < rr
	})
	return dto.NewLevelLayeredDiffTableHeaderDto(header.Entity(), sortedLevels, levelLayeredContent), nil
}

// Query specific difficult table's one level data contents with player related field (e.g PlayCount, Lamp status...)
//
// Requirements:
//
//	Level & ID & RivalID should not be empty
func (s *DiffTableService) QueryDiffTableDataWithRival(filter *vo.DiffTableHeaderVo) ([]*dto.DiffTableDataDto, int, error) {
	if filter.Level == "" {
		return nil, 0, fmt.Errorf("Level should not be empty")
	}
	if filter.ID <= 0 {
		return nil, 0, fmt.Errorf("ID should > 0")
	}
	if filter.RivalID <= 0 {
		return nil, 0, fmt.Errorf("RivalID should > 0")
	}
	rawContents, _, err := findDiffTableDataList(s.db, &vo.DiffTableDataVo{
		HeaderID:   filter.ID,
		Level:      filter.Level,
		Pagination: filter.Pagination,
	})
	if err != nil {
		return nil, 0, err
	}
	contents := dto.NewDiffTableDataDtoArray(rawContents)
	// NOTE: Here's a small hack to set correct play count in final result
	// "Play Count" column is setted by calling "mergeRivalRelatedData" method,
	// therefore if we merge "ghost"'s data first and "main user"'s data second,
	// "Play Count" column's data is always "main user"'s, not "ghost"'s.
	if filter.GhostRivalID > 0 {
		queryGhostScoreLogParam := &vo.RivalScoreLogVo{
			RivalId: filter.GhostRivalID,
		}
		if filter.GhostRivalTagID > 0 {
			tag, err := findRivalTagByID(s.db, filter.GhostRivalTagID)
			if err != nil {
				return nil, 0, err
			}
			queryGhostScoreLogParam.EndRecordTime = tag.RecordTime
		}
		ghostScoreLogsMap, err := findRivalScoreLogSha256Map(s.db, queryGhostScoreLogParam)
		if err != nil {
			return nil, 0, err
		}
		if err := mergeRivalRelatedData(ghostScoreLogsMap, contents, true); err != nil {
			return nil, 0, err
		}
	}
	sha256ScoreLogsMap, err := findRivalScoreLogSha256Map(s.db, &vo.RivalScoreLogVo{
		RivalId: filter.RivalID,
	})
	if err != nil {
		return nil, 0, err
	}
	if err := mergeRivalRelatedData(sha256ScoreLogsMap, contents, false); err != nil {
		return nil, 0, err
	}
	return contents, len(contents), nil
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

// Query if there exists a header that satisfies the condition
func queryDiffTableHeaderExistence(tx *gorm.DB, filter *entity.DiffTableHeader) (bool, error) {
	var dupCount int64
	if err := tx.Model(&entity.DiffTableHeader{}).Where(filter).Count(&dupCount).Error; err != nil {
		return false, err
	}
	return dupCount > 0, nil
}

func fetchDiffTableFromURL(url string) (*vo.DiffTableHeaderVo, error) {
	jsonUrl := ""
	splitUrl := strings.Split(url, "/")
	splitUrl[len(splitUrl)-1] = ""
	prefixUrl := strings.Join(splitUrl, "/")
	if strings.HasSuffix(url, ".html") {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return nil, err
			}
			line := strings.Trim(scanner.Text(), " ")
			// TODO: Any other cases?
			// Its pattern should be <meta name="bmstable" content="xxx.json" />
			if strings.HasPrefix(line, "<meta name=\"bmstable\"") {
				startp := strings.Index(line, "content") + len("content=\"") - 1
				if startp == -1 {
					log.Fatalf("Cannot find 'content' field in %s", url)
				}
				endp := -1
				// Finds the end position
				first := false
				for i := startp; i < len(line); i++ {
					if line[i] == '"' {
						if !first {
							first = true
						} else {
							endp = i
							break
						}
					}
				}
				if endp == -1 {
					log.Fatalf("Cannot find 'content' field in %s", url)
				}

				// Construct the json url path
				jsonUrl = prefixUrl + "/" + line[startp+1:endp]
				log.Debugf("Construct json url [%s] from [%s]", jsonUrl, url)
				break
			}
		}
	} else if strings.HasSuffix(url, ".json") {
		// Okay dokey
		jsonUrl = url
	}
	if jsonUrl == "" {
		return nil, fmt.Errorf("cannot fetch from %s", url)
	}
	var dth vo.DiffTableHeaderVo
	log.Debugf("before calling fetchJson, url=%s", jsonUrl)
	fetchJson(jsonUrl, &dth)
	if !strings.HasPrefix(dth.DataUrl, "http") {
		dth.DataUrl = prefixUrl + "/" + dth.DataUrl
	}
	return &dth, nil
}

func queryDiffTableInfoByID(tx *gorm.DB, ID uint) (*entity.DiffTableHeader, error) {
	var header entity.DiffTableHeader
	if err := tx.First(&header, ID).Error; err != nil {
		return nil, err
	}
	return &header, nil
}

// Merge player related data onto DiffTableDataDto (e.g PlayCount LampStatus...)
// TODO: We can actaully combine "query rival's related data" and "merge rival's data with DiffTableDataDto" two steps together
// The impede is mainly FindDiffTableHeaderListWithRival function, which requires redesign the data loading sequence
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
		if err := tx.Find(&headers).Error; err != nil {
			log.Error("[DiffTableService] Find difftable header failed with %v", err)
			return nil, 0, err
		}
		return headers, len(headers), nil
	}

	var headers []*entity.DiffTableHeader
	if err := tx.Where(filter.Entity()).Find(&headers).Error; err != nil {
		return nil, 0, err
	}
	return headers, len(headers), nil
}

// Query multiple difficult table's related level list by header ids
// When only related level list are required, this function is cheapier than load whole data content
//
// NOTE: parameter IDs must not be empty or the sql structure isn't correct
// Returns a list of pair(header_id, level)
func queryRelatedLevelByIDS(tx *gorm.DB, IDs []uint) (ret []struct {
	header_id uint
	level     string
}, err error) {
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
	sort.Slice(sorted, func(i, j int) bool {
		ll := sorted[i].Level
		rr := sorted[j].Level
		inxL := -1
		inxR := -1
		if preSortLevels != nil {
			inxL = slices.Index(preSortLevels, ll)
			inxR = slices.Index(preSortLevels, rr)
		}
		if inxL == -1 || inxR == -1 {
			ill, errL := strconv.Atoi(ll)
			irr, errR := strconv.Atoi(rr)
			if errL == nil && errR == nil {
				return ill < irr
			}
			return ll < rr
		}
		return inxL < inxR
	})
	return sorted
}

func fetchJson(url string, v interface{}) error {
	log.Debugf("Fetching json from url: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	// Hack \ufeff out, specially for PMS table
	body = bytes.ReplaceAll(body, []byte("\ufeff"), []byte(""))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}
