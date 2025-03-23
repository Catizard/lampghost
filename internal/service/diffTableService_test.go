package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Should be working on all real world tables
func TestFechDiffTableFromRealURL(t *testing.T) {
	var tests = []struct {
		name       string
		url        string
		symbol     string
		hasCourses bool
	}{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header, err := fetchDiffTableFromURL(tt.url)
			require.Nil(t, err)
			assert.Equal(t, header.Symbol, tt.symbol)
			err = header.ParseRawCourses()
			require.Nil(t, err)
			if tt.hasCourses && len(header.Courses) == 0 {
				t.Error("should have courses")
			}
		})
	}
}
