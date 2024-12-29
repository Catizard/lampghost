package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	if err := s.checkDuplicateHeaderUrl(url); err != nil {
		return nil, err
	}
	headerVo, err := fetchDiffTableFromURL(url)
	if err != nil {
		return nil, err
	}
	headerVo.HeaderUrl = url
	if headerVo.DataUrl == "" {
		return nil, fmt.Errorf("assert: header.DataUrl cannot be empty")
	}
	log.Debugf("[DiffTableService] Got header data: %v", headerVo)
	if err := s.checkDuplicateDataUrl(headerVo.DataUrl); err != nil {
		return nil, err
	}

	// Transaction begins from here
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// (1) difficult table header
	headerEntity := headerVo.Entity()
	if err := tx.Create(headerEntity).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// (2) difficult related course contents
	if len(headerVo.Courses) > 0 {
		var courseData []entity.CourseInfo
		for _, arr := range headerVo.Courses {
			for _, courseInfoVo := range arr {
				courseInfo := courseInfoVo.Entity()
				courseInfo.HeaderID = headerEntity.ID
				courseData = append(courseData, *courseInfo)
			}
		}
		if err := tx.Unscoped().Where("header_id = ?", headerEntity.ID).Delete(&entity.CourseInfo{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Create(&courseData).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	// (3) difficult table concreate contents
	var data []entity.DiffTableData
	if err := fetchJson(headerVo.DataUrl, &data); err != nil {
		return nil, err
	}
	for i := range data {
		data[i].HeaderID = headerEntity.ID
	}
	if err := tx.Unscoped().Where("header_id = ?", headerEntity.ID).Delete(&entity.DiffTableData{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Transaction ends here
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	log.Infof("[DiffTableService] Inserted one header with %d contents", len(data))
	return headerEntity, nil
}

// Query all difficult table datas
//
// Returns difficult header and its contents
func (s *DiffTableService) FindDiffTableHeaderList() ([]dto.DiffTableHeaderDto, int, error) {
	var headers []entity.DiffTableHeader
	if err := s.db.Find(&headers).Error; err != nil {
		log.Error("[DiffTableService] Find difftable header failed with %v", err)
		return nil, 0, err
	}
	headerIds := make([]uint, len(headers))
	for i, header := range headers {
		headerIds[i] = header.ID
	}
	rawContents, err := queryDiffTableDataByIDs(s.db, headerIds)
	if err != nil {
		return nil, 0, err
	}

	ret := make([]dto.DiffTableHeaderDto, len(headers))
	for i, header := range headers {
		contents := make([]dto.DiffTableDataDto, 0)
		for _, content := range rawContents {
			if content.HeaderID == header.ID {
				contents = append(contents, content)
			}
		}
		ret[i] = *dto.NewDiffTableHeaderDto(&header, contents)
	}

	return ret, len(ret), nil
}

// Extend function for FindDiffTableHeaderList
//
// Adds player related field (e.g PlayCount, Lamp status)
func (s *DiffTableService) FindDiffTableHeaderListWithRival(rivalID uint) ([]dto.DiffTableHeaderDto, int, error) {
	headers, _, err := s.FindDiffTableHeaderList()
	if err != nil {
		return nil, 0, err
	}
	songLogs, _, err := findRivalScoreLogList(s.db, rivalID)
	if err != nil {
		return nil, 0, err
	}
	sha256ScoreLogsMap := make(map[string][]entity.RivalScoreLog)
	for _, scoreLog := range songLogs {
		if _, ok := sha256ScoreLogsMap[scoreLog.Sha256]; !ok {
			sha256ScoreLogsMap[scoreLog.Sha256] = make([]entity.RivalScoreLog, 0)
		}
		sha256ScoreLogsMap[scoreLog.Sha256] = append(sha256ScoreLogsMap[scoreLog.Sha256], *scoreLog)
	}
	for i, header := range headers {
		for j, content := range header.Contents {
			if logs, ok := sha256ScoreLogsMap[content.Sha256]; ok {
				headers[i].Contents[j].PlayCount = len(logs)
				for _, log := range logs {
					headers[i].Contents[j].Lamp = max(content.Lamp, int(log.Clear))
				}
			}
		}
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
	contents, err := queryDiffTableDataByHeaderID(s.db, ID)
	if err != nil {
		return nil, err
	}
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
	songLogs, _, err := findRivalScoreLogList(s.db, rivalID)
	if err != nil {
		return nil, err
	}
	sha256ScoreLogsMap := make(map[string][]entity.RivalScoreLog)
	for _, scoreLog := range songLogs {
		if _, ok := sha256ScoreLogsMap[scoreLog.Sha256]; !ok {
			sha256ScoreLogsMap[scoreLog.Sha256] = make([]entity.RivalScoreLog, 0)
		}
		sha256ScoreLogsMap[scoreLog.Sha256] = append(sha256ScoreLogsMap[scoreLog.Sha256], *scoreLog)
	}
	for i, content := range header.Contents {
		if logs, ok := sha256ScoreLogsMap[content.Sha256]; ok {
			header.Contents[i].PlayCount = len(logs)
			for _, log := range logs {
				header.Contents[i].Lamp = max(content.Lamp, int(log.Clear))
			}
		}
	}
	return header, nil
}

func (s *DiffTableService) QueryLevelLayeredDiffTableInfoById(ID uint) (*dto.DiffTableHeaderDto, error) {
	header, err := s.QueryDiffTableInfoByID(ID)
	if err != nil {
		return nil, err
	}
	levels := make(map[string]interface{})
	levelLayeredContent := make(map[string][]dto.DiffTableDataDto)
	for _, v := range header.Contents {
		if _, ok := levelLayeredContent[v.Level]; !ok {
			levelLayeredContent[v.Level] = make([]dto.DiffTableDataDto, 0)
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

func (s *DiffTableService) checkDuplicateHeaderUrl(headerUrl string) error {
	var dupCount int64
	if err := s.db.Model(&entity.DiffTableHeader{}).Where("header_url = ?", headerUrl).Count(&dupCount).Error; err != nil {
		return err
	}
	if dupCount > 0 {
		return fmt.Errorf("header url: %s is duplicated", headerUrl)
	}
	return nil
}

func (s *DiffTableService) checkDuplicateDataUrl(dataUrl string) error {
	var dupCount int64
	if err := s.db.Model(&entity.DiffTableHeader{}).Where("data_url = ?", dataUrl).Count(&dupCount).Error; err != nil {
		return err
	}
	if dupCount > 0 {
		return fmt.Errorf("data url: %s is duplicated", dataUrl)
	}
	return nil
}

func fetchDiffTableFromURL(url string) (*vo.DiffTableHeaderVo, error) {
	jsonUrl := ""
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
				splitUrl := strings.Split(url, "/")
				splitUrl[len(splitUrl)-1] = line[startp+1 : endp]
				jsonUrl = strings.Join(splitUrl, "/")
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
	return &dth, nil
}

func queryDiffTableInfoByID(tx *gorm.DB, ID uint) (*entity.DiffTableHeader, error) {
	var header entity.DiffTableHeader
	if err := tx.First(&header, ID).Error; err != nil {
		return nil, err
	}
	log.Debugf("[DiffTableService] QueryDiffTableInfoByID fetched header: %v", header)
	return &header, nil
}

// Query specific difficult table's all contents
// Note this function directly returns dto instead of entity form of data, which is an
// old problem that hash field is incompitable in some difficult tables
//
// This function uses default user's song data to build the cache
func queryDiffTableDataByHeaderID(tx *gorm.DB, headerID uint) ([]dto.DiffTableDataDto, error) {
	var rawContents []entity.DiffTableData
	if err := tx.Where("header_id = ?", headerID).Find(&rawContents).Error; err != nil {
		return nil, err
	}
	cache, err := queryDefaultSongHashCache(tx)
	if err != nil {
		return nil, err
	}
	contents := make([]dto.DiffTableDataDto, 0)
	for _, rawContent := range rawContents {
		contents = append(contents, *dto.NewDiffTableDataDtoWithCache(&rawContent, cache))
	}
	return contents, nil
}

// Query multiple difficult table's contents by header ids
//
// Extends to queryDiffTableDataByHeaderID, which could query multiple ids
func queryDiffTableDataByIDs(tx *gorm.DB, IDs []uint) ([]dto.DiffTableDataDto, error) {
	var rawContents []entity.DiffTableData
	if err := tx.Debug().Where("header_id in ?", IDs).Find(&rawContents).Error; err != nil {
		return nil, err
	}
	cache, err := queryDefaultSongHashCache(tx)
	if err != nil {
		return nil, err
	}
	contents := make([]dto.DiffTableDataDto, 0)
	for _, rawContent := range rawContents {
		contents = append(contents, *dto.NewDiffTableDataDtoWithCache(&rawContent, cache))
	}
	return contents, nil
}

func fetchJson(url string, v interface{}) error {
	log.Debugf("Fetching json from url: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}
