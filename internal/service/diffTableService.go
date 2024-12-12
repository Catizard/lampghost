package service

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
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

func (s *DiffTableService) AddDiffTableHeader(url string) (*entity.DiffTableHeader, error) {
	url = strings.TrimSpace(url)
	header, err := fetchDiffTableFromURL(url)
	if err != nil {
		return nil, err
	}
	// TODO: handle duplication
	var data []entity.DiffTableData
	if err := fetchJson(header.DataUrl, data); err != nil {
		return nil, err
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(header).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("header_id = ?", header.ID).Delete(&entity.DiffTableData{}).Error; err != nil {
			return err
		}
		if err := tx.Create(&data).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return header, nil
}

func fetchDiffTableFromURL(url string) (*entity.DiffTableHeader, error) {
	jsonUrl := ""
	if strings.HasSuffix(url, ".html") {
		log.Infof("Fetch difficult table data from %s", url)
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
		log.Fatalf("Cannot fetch %s", url)
	}
	dth := &entity.DiffTableHeader{}
	fetchJson(jsonUrl, &dth)
	return dth, nil
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
