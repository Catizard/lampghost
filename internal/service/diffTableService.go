package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type DiffTableService struct {
	db *gorm.DB
	rivalSongDataService *RivalSongDataService
}

func NewDiffTableService(db *gorm.DB, rivalSongDataService *RivalSongDataService) *DiffTableService {
	return &DiffTableService{
		db: db,
		rivalSongDataService: rivalSongDataService,
	}
}

func (s *DiffTableService) AddDiffTableHeader(url string) (*entity.DiffTableHeader, error) {
	url = strings.TrimSpace(url)
	log.Debugf("[DiffTableService] calling AddDiffTableHeader with url: %s", url)
	if err := s.checkDuplicateHeaderUrl(url); err != nil {
		return nil, err
	}
	header, err := fetchDiffTableFromURL(url)
	if err != nil {
		return nil, err
	}
	header.HeaderUrl = url
	if err != nil {
		return nil, err
	}
	if header.DataUrl == "" {
		return nil, fmt.Errorf("assert: header.DataUrl cannot be empty")
	}
	log.Debugf("[DiffTableService] Got header data: %s", header)
	if err := s.checkDuplicateDataUrl(header.DataUrl); err != nil {
		return nil, err
	}
	var data []entity.DiffTableData
	if err := fetchJson(header.DataUrl, &data); err != nil {
		return nil, err
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(header).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("header_id = ?", header.ID).Delete(&entity.DiffTableData{}).Error; err != nil {
			return err
		}
		for i := range data {
			data[i].HeaderID = header.ID
		}
		// TODO: Create course info here
		if err := tx.Create(&data).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Errorf("[DiffTableService] Add difftable header failed with %v", err)
		return nil, err
	}
	log.Infof("[DiffTableService] Inserted one header with %d contents", len(data))
	return header, nil
}

func (s *DiffTableService) FindDiffTableHeaderList() ([]dto.DiffTableHeaderDto, int, error) {
	var headers []entity.DiffTableHeader
	if err := s.db.Find(&headers).Error; err != nil {
		log.Error("[DiffTableService] Find difftable header failed with %v", err)
		return nil, 0, err
	}
	ret := make([]dto.DiffTableHeaderDto, len(headers))
	for i, header := range headers {
		ret[i] = *dto.NewDiffTableHeaderDto(&header, nil)
	}
	return ret, len(ret), nil
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

func (s *DiffTableService) QueryDiffTableInfoByID(ID uint) (*dto.DiffTableHeaderDto, error) {
	var header entity.DiffTableHeader
	if err := s.db.Find(&header, ID).Error; err != nil {
		return nil, err
	}
	var rawContents []entity.DiffTableData
	if err := s.db.Where(&entity.DiffTableData{HeaderID: ID}).Find(&rawContents).Error; err != nil {
		return nil, err
	}
	cache, err := s.rivalSongDataService.QueryDefaultSongHashCache()
	if err != nil {
		return nil, err
	}
	// NOTE: 如果一个数据找不到对应的hash信息，则认为剔除掉该数据避免后续出错
	contents := make([]entity.DiffTableData, 0)
	for _, rawContent := range rawContents {
		if rawContent.Sha256 != "" {
			if md5, ok := cache.GetMD5(rawContent.Sha256); ok {
				rawContent.Md5 = md5
				contents = append(contents, rawContent)
			}
		} else if rawContent.Md5 != "" {
			if sha256, ok := cache.GetSHA256(rawContent.Md5); ok { 
				rawContent.Sha256 = sha256
				contents = append(contents, rawContent)
			}
		}
	}

	return dto.NewDiffTableHeaderDto(&header, contents), nil
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

func fetchDiffTableFromURL(url string) (*entity.DiffTableHeader, error) {
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
	dth := entity.DiffTableHeader{}
	log.Debugf("before calling fetchJson, url=%s", jsonUrl)
	fetchJson(jsonUrl, &dth)
	return &dth, nil
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
