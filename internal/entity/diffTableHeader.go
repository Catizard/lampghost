package entity

import (
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

	headerUrls []string // Internal use, to simplify the init step
}

// One predefined header scheme is basically an array of headers
// and some additional meta info
type PredefineTableScheme struct {
	Headers []PredefineTableHeader
	Name    string
	//SelectedCategory []string
}

var PredefineTableSchemeNames []string = []string{
	"Raw",
	"Zris",
}

var PredefineTableSchemes map[string]PredefineTableScheme = make(map[string]PredefineTableScheme)

func init() {
	urlJoinedHeaders := make([]PredefineTableHeader, 0)
	// NOTE: The element in headerUrls' order must follow the order of PredefineTableSchemeNames
	// Starter category
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "Stardust",
		Symbol:       "ξ",
		Category:     "Starter",
		TagColor:     "#E5F4D3",
		TagTextColor: "#2B410D",
		headerUrls: []string{
			"https://mqppppp.neocities.org/StardustTable.html",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "Starlight",
		Symbol:       "sr",
		Category:     "Starter",
		TagColor:     "#E8EDFF",
		TagTextColor: "#121D52",
		headerUrls: []string{
			"https://djkuroakari.github.io/starlighttable.html",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "通常難易度表",
		Symbol:       "☆",
		Category:     "Starter",
		TagColor:     "#EDF7FD",
		TagTextColor: "#0F5F8A",
		headerUrls: []string{
			"http://www.ribbit.xyz/bms/tables/normal.html",
			"http://zris.work/bmstable/normal/normal_header.json",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "NEW GENERATION 通常難易度表",
		Symbol:       "▽",
		Category:     "Starter",
		TagColor:     "#FFFAEB",
		TagTextColor: "#946D18",
		headerUrls: []string{
			"http://rattoto10.jounin.jp/table.html",
			"http://zris.work/bmstable/normal2/header.json",
		},
	})
	// Insane category
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "Satellite",
		Symbol:       "sl",
		Category:     "Insane",
		TagColor:     "#B6EAD2",
		TagTextColor: "#0A4D2F",
		headerUrls: []string{
			"https://stellabms.xyz/sl/table.html",
			"http://zris.work/bmstable/satellite/header.json",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "発狂BMS難易度表",
		Symbol:       "★",
		Category:     "Insane",
		TagColor:     "#A3D8F5",
		TagTextColor: "#0F5F8A",
		headerUrls: []string{
			"http://www.ribbit.xyz/bms/tables/insane.html",
			"http://zris.work/bmstable/insane/insane_header.json",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "NEW GENERATION 発狂難易度表",
		Symbol:       "▼",
		Category:     "Insane",
		TagColor:     "#FFECB9",
		TagTextColor: "#946D18",
		headerUrls: []string{
			"http://rattoto10.jounin.jp/table_insane.html",
			"http://zris.work/bmstable/insane2/insane_header.json",
		},
	})
	// Overjoy category
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "第三期Overjoy",
		Symbol:       "★★",
		Category:     "Overjoy",
		TagColor:     "#DDBDF2",
		TagTextColor: "#5C2989",
		headerUrls: []string{
			"http://rattoto10.jounin.jp/table_overjoy.html",
			"http://zris.work/bmstable/overjoy/header.json",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "Stella",
		Symbol:       "st",
		Category:     "Overjoy",
		TagColor:     "#FFB5A8",
		TagTextColor: "#331510",
		headerUrls: []string{
			"https://stellabms.xyz/st/table.html",
			"http://zris.work/bmstable/stella/header.json",
		},
	})
	// DP category
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "δ難易度表",
		Symbol:       "δ",
		Category:     "DP",
		TagColor:     "#EDF7FD",
		TagTextColor: "#0F5F8A",
		headerUrls: []string{
			"https://deltabms.yaruki0.net/table/dpdelta",
			"http://zris.work/bmstable/dp_normal/dpn_header.json",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "発狂DP難易度表",
		Symbol:       "★",
		Category:     "DP",
		TagColor:     "#A3D8F5",
		TagTextColor: "#0F5F8A",
		headerUrls: []string{
			"https://deltabms.yaruki0.net/table/dpinsane",
			"http://zris.work/bmstable/dp_insane/dpi_header.json",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "DP Overjoy",
		Symbol:       "★★",
		Category:     "DP",
		TagColor:     "#DDBDF2",
		TagTextColor: "#5C2989",
		headerUrls: []string{
			"http://ereter.net/dpoverjoy",
			"http://zris.work/bmstable/dp_overjoy/header.json",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "DP Satellite",
		Symbol:       "DPsl",
		Category:     "DP",
		TagColor:     "#B6EAD2",
		TagTextColor: "#0A4D2F",
		headerUrls: []string{
			"https://stellabms.xyz/dp/table.html",
			"http://zris.work/bmstable/dp_satellite/header.json",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "DP Stella",
		Symbol:       "DPst",
		Category:     "DP",
		TagColor:     "#FFB5A8",
		TagTextColor: "#331510",
		headerUrls: []string{
			"https://stellabms.xyz/dpst/table.html",
			"http://zris.work/bmstable/dp_stella/header.json",
		},
	})
	// PMS category
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "PMSデータベース(Lv1~45)",
		Symbol:       "PLv",
		Category:     "PMS",
		TagColor:     "#EDF7FD",
		TagTextColor: "#0F5F8A",
		headerUrls: []string{
			"http://pmsdifficulty.xxxxxxxx.jp/PMSdifficulty.html",
			"http://zris.work/bmstable/pms_normal/pmsdatabase_header.json",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "発狂PMSデータベース(lv46～)",
		Symbol:       "P●",
		Category:     "PMS",
		TagColor:     "#A3D8F5",
		TagTextColor: "#0F5F8A",
		headerUrls: []string{
			"https://pmsdifficulty.xxxxxxxx.jp/insane_PMSdifficulty.html",
			"http://zris.work/bmstable/pms_insane/insane_pmsdatabase_header.json",
		},
	})
	urlJoinedHeaders = append(urlJoinedHeaders, PredefineTableHeader{
		Name:         "発狂PMS難易度表",
		Symbol:       "●",
		Category:     "PMS",
		TagColor:     "#FFECB9",
		TagTextColor: "#946D18",
		headerUrls: []string{
			"https://pmsdifficulty.xxxxxxxx.jp/_pastoral_home.html",
			"http://zris.work/bmstable/pms_upper/header.json",
		},
	})

	for i, schemeName := range PredefineTableSchemeNames {
		headers := make([]PredefineTableHeader, 0)
		for _, header := range urlJoinedHeaders {
			shallow := header
			// push down
			if i < len(header.headerUrls) {
				shallow.HeaderUrl = header.headerUrls[i]
			} else {
				// If not provided, fallback to use the raw one
				shallow.HeaderUrl = header.headerUrls[0]
			}
			headers = append(headers, shallow)
		}
		PredefineTableSchemes[schemeName] = PredefineTableScheme{
			Name:    schemeName,
			Headers: headers,
		}
	}
}
