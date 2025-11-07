package entity

import (
    "encoding/json"
    "strings"
    "embed"

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
    //SelectedCategory []string
}

var PredefineTableSchemes map[string]PredefineTableScheme = make(map[string]PredefineTableScheme)

//go:embed predefine_table_schemes.json
var predefineTableSchemesFS embed.FS

type predefineTableSchemesFile struct {
    SchemeNames []string `json:"schemeNames"`
    Headers     []struct {
        Name         string            `json:"name"`
        Symbol       string            `json:"symbol"`
        Category     string            `json:"category"`
        TagColor     string            `json:"tagColor"`
        TagTextColor string            `json:"tagTextColor"`
        HeaderUrls   map[string]string `json:"headerUrls"`
    } `json:"headers"`
}

func init() {
    var file predefineTableSchemesFile
    data, err := predefineTableSchemesFS.ReadFile("predefine_table_schemes.json")
    if err != nil {
        panic(err)
    }
    if err := json.Unmarshal(data, &file); err != nil {
        panic(err)
    }
    for _, schemeName := range file.SchemeNames {
        headers := make([]PredefineTableHeader, 0, len(file.Headers))
        for _, h := range file.Headers {
            shallow := PredefineTableHeader{
                Name:         h.Name,
                Symbol:       h.Symbol,
                Category:     h.Category,
                TagColor:     h.TagColor,
                TagTextColor: h.TagTextColor,
            }
            url := h.HeaderUrls["Raw"]
            if u, ok := h.HeaderUrls[schemeName]; ok && u != "" {
                url = u
            }
            shallow.HeaderUrl = url
            headers = append(headers, shallow)
        }
        PredefineTableSchemes[schemeName] = PredefineTableScheme{
            Name:    schemeName,
            Headers: headers,
        }
    }
}
