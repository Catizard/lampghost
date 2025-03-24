package service

import (
	"testing"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type EmptyMemoryDatabaseTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *EmptyMemoryDatabaseTestSuite) SetupTest() {
	dbConfig := config.DatabaseConfig{
		DSN: ":memory:",
	}
	db, err := database.NewDatabase(&dbConfig)
	suite.NoError(err)
	suite.db = db
}

func (s *EmptyMemoryDatabaseTestSuite) TearDownTest() {
}

type diffTableDefinition struct {
	name       string
	url        string
	symbol     string
	hasCourses bool
}

// Real table definition
var realTableDefintion = [...]diffTableDefinition{
	{"通常難易度表", "http://zris.work/bmstable/normal/normal_header.json", "☆", false},
	{"発狂BMS難易度表", "http://zris.work/bmstable/insane/insane_header.json", "★", false},
	{"第三期Overjoy", "http://zris.work/bmstable/overjoy/header.json", "★★", false},
	{"NEW GENERATION 通常難易度表", "http://zris.work/bmstable/normal2/header.json", "▽", true},
	{"NEW GENERATION 発狂難易度表", "http://zris.work/bmstable/insane2/insane_header.json", "▼", true},
	{"Satellite", "http://zris.work/bmstable/satellite/header.json", "sl", true},
	{"Stella", "http://zris.work/bmstable/stella/header.json", "st", true},
	{"DP Satellite", "http://zris.work/bmstable/dp_satellite/header.json", "DPsl", true},
	{"DP Stella", "http://zris.work/bmstable/dp_stella/header.json", "DPst", false},
	{"δ難易度表", "http://zris.work/bmstable/dp_normal/dpn_header.json", "δ", true},
	{"発狂DP難易度表", "http://zris.work/bmstable/dp_insane/dpi_header.json", "★", true},
	{"DP Overjoy", "http://zris.work/bmstable/dp_overjoy/header.json", "★★", false},
	{"DPBMS白難易度表(通常)", "http://zris.work/bmstable/dp_white/header.json", "白", false},
	{"DPBMS黒難易度表(発狂)", "http://zris.work/bmstable/dp_black/header.json", "黒", false},
	{"発狂DPBMSごった煮難易度表", "http://zris.work/bmstable/dp_zhu/header.json", "★", false},
	{"発狂14keyBMS闇鍋難易度表", "http://zris.work/bmstable/dp_anguo/head14.json", "★", false},
	{"DPBMSと諸感", "http://zris.work/bmstable/dp_zhugan/header.json", "☆", false},
	{"Luminous", "http://zris.work/bmstable/luminous/header.json", "ln", false},
	{"LN難易度", "http://zris.work/bmstable/ln/ln_header.json", "◆", true},
	{"Scramble難易度表", "http://zris.work/bmstable/scramble/header.json", "SB", true},
	{"PMSデータベース(Lv1~45)", "http://zris.work/bmstable/pms_normal/pmsdatabase_header.json", "PLv", false},
	{"発狂PMSデータベース(lv46～)", "http://zris.work/bmstable/pms_insane/insane_pmsdatabase_header.json", "P●", false},
	{"発狂PMS難易度表", "http://zris.work/bmstable/pms_upper/header.json", "●", true},
	// This table has been lost
	// {"PMS Database コースデータ案内所", "http://zris.work/bmstable/pms_course/course_header.jsonn", "Pcourse", false},
	{"Stellalite", "http://zris.work/bmstable/stellalite/Stellalite-header.json", "stl", false},
	{"オマージュBMS難易度表", "http://zris.work/bmstable/homage/header.json", "∽", false},
}

// Should be working on all real world tables
func TestFetchDiffTableFromRealURL(t *testing.T) {
	for _, tt := range realTableDefintion {
		t.Run(tt.name, func(t *testing.T) {
			header, err := fetchDiffTableFromURL(tt.url)
			require.Nil(t, err)
			assert.Equal(t, header.Symbol, tt.symbol)
			err = header.ParseRawCourses()
			require.Nil(t, err)
			if tt.hasCourses {
				require.NotEqual(t, 0, len(header.Courses))
				for _, course := range header.Courses {
					assert.NotEmpty(t, course.Md5)
				}
			}
		})
	}
}

// Basic test for `AddDiffTableHeader` method
// This test is kind of overlapping `TestFetchDiffTableFromRealURL`, so no all real url is being considered as test cases
// Only a few non-standard or buggy url are invovled
func (suite *EmptyMemoryDatabaseTestSuite) TestAddDiffTableHeader() {
	tests := []struct {
		name       string
		url        string
		hasCourses bool
	}{
		{"NEW GENERATION 発狂難易度表", "http://zris.work/bmstable/normal2/header.json", true},
		{"δ難易度表", "http://zris.work/bmstable/dp_normal/dpn_header.json", true},
		{"発狂DP難易度表", "http://zris.work/bmstable/dp_insane/dpi_header.json", true},
	}
	diffTableService := NewDiffTableService(suite.db)
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			header, err := diffTableService.AddDiffTableHeader(tt.url)
			require.Nil(t, err)
			assert.NotEqual(t, 0, header.ID)
		})
		// NOTE: we cannot use CourseInfoService here to query inserted courses, because `FindCourseInfoList` method requires data from `RivalSongData`.
	}
}

func TestDatabaseRelatedWrapper(t *testing.T) {
	suite.Run(t, new(EmptyMemoryDatabaseTestSuite))
}
