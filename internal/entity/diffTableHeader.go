package entity

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/Catizard/bmstable"
	"gorm.io/gorm"
)

type DiffTableHeader struct {
	gorm.Model

	HeaderUrl    string
	DataUrl      string
	Name         string
	OriginalUrl  *string
	Symbol       string
	OrderNumber  int `gorm:"default:0"`
	LevelOrders  string
	TagColor     string
	TagTextColor string
	NoTagBuild   *int `gorm:"default:0"`
}

func (DiffTableHeader) TableName() string {
	return "difftable_header"
}

// Convert external difficult table definition to internal one
// If inheritHeader is non-nil, inherit some extra fields from it (esp color definitions)
func NewDiffTableHeaderFromImport(importHeader *bmstable.DifficultTable, inheritHeader *DiffTableHeader) *DiffTableHeader {
	ret := &DiffTableHeader{
		HeaderUrl:   importHeader.HeaderURL,
		DataUrl:     importHeader.DataURL,
		Name:        importHeader.Name,
		OriginalUrl: &importHeader.OriginalURL,
		Symbol:      importHeader.Symbol,
		LevelOrders: strings.Join(importHeader.LevelOrder, ","),
	}
	ret.TagColor = inheritHeader.TagColor
	ret.TagTextColor = inheritHeader.TagTextColor
	ret.NoTagBuild = inheritHeader.NoTagBuild
	return ret
}

// Cropped struct definition of DiffTableHeader
type PredefineTableHeader struct {
	HeaderUrl    string
	Name         string
	Symbol       string
	TagColor     string
	TagTextColor string
	Category     string
}

// One predefined header scheme is basically an array of headers
// and some additional meta info
type PredefineTableScheme struct {
	Headers []PredefineTableHeader
	Name    string
}

//go:embed data/predefine_tables.json
var predefineTablesJSON []byte

type predefineTableJSON struct {
	SchemeNames []string                 `json:"schemeNames"`
	Headers     []predefineTableJSONEntry `json:"headers"`
}

type predefineTableJSONEntry struct {
	Name         string            `json:"name"`
	Symbol       string            `json:"symbol"`
	Category     string            `json:"category"`
	TagColor     string            `json:"tagColor"`
	TagTextColor string            `json:"tagTextColor"`
	HeaderUrls   map[string]string `json:"headerUrls"`
}

var PredefineTableSchemeNames []string
var PredefineTableSchemes map[string]PredefineTableScheme

func init() {
	var jsonData predefineTableJSON
	if err := json.Unmarshal(predefineTablesJSON, &jsonData); err != nil {
		panic("failed to parse predefine_tables.json: " + err.Error())
	}

	PredefineTableSchemeNames = jsonData.SchemeNames
	PredefineTableSchemes = make(map[string]PredefineTableScheme, len(PredefineTableSchemeNames))

	for _, schemeName := range PredefineTableSchemeNames {
		headers := make([]PredefineTableHeader, 0, len(jsonData.Headers))
		for _, h := range jsonData.Headers {
			header := PredefineTableHeader{
				Name:         h.Name,
				Symbol:       h.Symbol,
				Category:     h.Category,
				TagColor:     h.TagColor,
				TagTextColor: h.TagTextColor,
				HeaderUrl:    h.HeaderUrls["Raw"],
			}
			if url, ok := h.HeaderUrls[schemeName]; ok {
				header.HeaderUrl = url
			}
			headers = append(headers, header)
		}
		PredefineTableSchemes[schemeName] = PredefineTableScheme{
			Name:    schemeName,
			Headers: headers,
		}
	}
}
