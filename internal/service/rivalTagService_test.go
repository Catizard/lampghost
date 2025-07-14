package service_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Some tables' courses are defined by using sha256 instead of md5
// This test ensures that they are compatible when building tags
func TestSha256CoursesTagBuild(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	db, err := database.NewMemoryDatabase()
	require.Nil(t, err)
	service := service.NewRivalTagService(db)
	// Mock a difficult table which has two courses, defined as below:
	// NOTE: for convient, we define sha256 == md5
	// Course-1: sha256s="1,2,3,4", md5s=""
	// Course-2: sha256s="2,3,4,5", md5s=""
	// Course-3: sha256s="", md5s="1,2,3,4"
	// Course-4: sha256s="", md5s="2,3,4,5"
	// NOTE: however, we cannot use "1,2,3,4", because the real sha256's length is 64, while the inner implementation is based on this
	// Therefore we need to extend them to 64 length too
	one_64 := "1111111111111111111111111111111111111111111111111111111111111111"
	two_64 := "2222222222222222222222222222222222222222222222222222222222222222"
	three_64 := "3333333333333333333333333333333333333333333333333333333333333333"
	four_64 := "4444444444444444444444444444444444444444444444444444444444444444"
	five_64 := "5555555555555555555555555555555555555555555555555555555555555555"

	header := &entity.DiffTableHeader{
		Name: "test",
	}
	require.Nil(t, db.Create(&header).Error)
	courseDatas := []entity.CourseInfo{
		{HeaderID: header.ID, Name: "Course-1", Sha256s: strings.Join([]string{one_64, two_64, three_64, four_64}, ","), Md5s: ""},
		{HeaderID: header.ID, Name: "Course-2", Sha256s: strings.Join([]string{two_64, three_64, four_64, five_64}, ","), Md5s: ""},
		{HeaderID: header.ID, Name: "Course-3", Sha256s: "", Md5s: strings.Join([]string{one_64, two_64, three_64, four_64}, ",")},
		{HeaderID: header.ID, Name: "Course-4", Sha256s: "", Md5s: strings.Join([]string{two_64, three_64, four_64, five_64}, ",")},
	}
	require.Nil(t, db.Create(&courseDatas).Error)
	// tags are based on rivals, therefore we need a default user and some logs for testing
	require.Nil(t, db.Create(&entity.RivalInfo{
		Name:     "-",
		MainUser: true,
	}).Error)
	require.Nil(t, db.Create(&entity.RivalScoreLog{
		RivalId:    1,
		Sha256:     strings.Join([]string{one_64, two_64, three_64, four_64}, ""),
		Clear:      entity.Hard,
		Mode:       "0",
		RecordTime: time.Now(),
	}).Error)
	t.Run("OneLog-NoSongData", func(t *testing.T) {
		err := service.SyncRivalTag(1)
		require.Nil(t, err)
		tags, n, err := service.FindRivalTagList(&vo.RivalTagVo{
			RivalId: 1,
		})
		require.Nil(t, err)
		require.Equal(t, 1, n)
		assert.Equal(t, tags[0].Generated, true)
		assert.Equal(t, tags[0].TagName, "Course-1 First Clear")
	})
	// fill the blank between sha256 to md5
	// 1 -> 1; 2 -> 2; 3 -> 3; 4 -> 4;
	songDatas := []entity.RivalSongData{
		{RivalId: 1, Sha256: one_64, Md5: one_64},
		{RivalId: 1, Sha256: two_64, Md5: two_64},
		{RivalId: 1, Sha256: three_64, Md5: three_64},
		{RivalId: 1, Sha256: four_64, Md5: four_64},
	}
	require.Nil(t, db.Create(&songDatas).Error)
	// TODO: Can we get rid of this?
	// // NOTE: we need expire the default cache after writing data into songdata table
	// expireDefaultCache()
	// t.Run("OneLog-HasSongData", func(t *testing.T) {
	// 	err := service.SyncRivalTag(1)
	// 	require.Nil(t, err)
	// 	_, n, err := service.FindRivalTagList(&vo.RivalTagVo{
	// 		RivalId: 1,
	// 	})
	// 	require.Nil(t, err)
	// 	require.Equal(t, 2, n)
	// })
}

func TestAddRivalTag(t *testing.T) {
	db, err := database.NewMemoryDatabase()
	if err != nil {
		t.Fatalf("db: %s", err)
	}
	rivalTagService := service.NewRivalTagService(db)
	rivalInfoService := newRivalInfoService(db)
	emptyMainUser := newEmptyInitializeUser(false, false, false)
	if err := rivalInfoService.InitializeMainUser(emptyMainUser); err != nil {
		t.Fatalf("initializeMainUser: %s", err)
	}

	t.Run("quick fails", func(t *testing.T) {
		tests := []struct {
			name  string
			input *vo.RivalTagVo
		}{
			{
				"completely nil",
				nil,
			},
			{
				"rival id is 0",
				&vo.RivalTagVo{RivalId: 0, TagName: "TEST", RecordTime: time.Now()},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := rivalTagService.AddRivalTag(tt.input); err == nil {
					t.Fatalf("expected error, got nothing")
				}
			})
		}
	})

	t.Run("smoke", func(t *testing.T) {
		if err := rivalTagService.AddRivalTag(&vo.RivalTagVo{
			RivalId:    1,
			TagName:    "TEST",
			RecordTime: time.Now(),
		}); err != nil {
			t.Fatalf("add rival tag: %s", err)
		}
	})
}

func TestUpdateRivalTag(t *testing.T) {
	db, err := database.NewMemoryDatabase()
	if err != nil {
		t.Fatalf("db: %s", err)
	}
	rivalTagService := service.NewRivalTagService(db)
	rivalInfoService := newRivalInfoService(db)
	emptyMainUser := newEmptyInitializeUser(false, false, false)
	if err := rivalInfoService.InitializeMainUser(emptyMainUser); err != nil {
		t.Fatalf("initializeMainUser: %s", err)
	}

	t.Run("quick fails", func(t *testing.T) {
		tests := []struct {
			name  string
			input *vo.RivalTagUpdateParam
		}{
			{"completely nil", nil},
			{"no id", &vo.RivalTagUpdateParam{}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := rivalTagService.UpdateRivalTag(tt.input); err == nil {
					t.Fatalf("expected error, got nothing")
				}
			})
		}
	})

	t.Run("smoke", func(t *testing.T) {
		currentTime := time.Now()
		if err := rivalTagService.AddRivalTag(&vo.RivalTagVo{
			RivalId:    1,
			TagName:    "TEST",
			RecordTime: currentTime,
		}); err != nil {
			t.Fatalf("add rival tag: %s", err)
		}

		tags, n, err := rivalTagService.FindRivalTagList(nil)
		if err != nil {
			t.Fatalf("FindRivalTagList: %s", err)
		}
		if n != 1 {
			t.Fatalf("expected exactly one rival tag, got %d", n)
		}
		tag := tags[0]
		t.Run("update", func(t *testing.T) {
			newName := "NEW NAME"
			newSymbol := "NEW SYMBOL"
			newEnabled := true
			if err := rivalTagService.UpdateRivalTag(&vo.RivalTagUpdateParam{
				ID:         tag.ID,
				TagName:    &newName,
				RecordTime: time.Now(),
				Symbol:     &newSymbol,
				Enabled:    &newEnabled,
			}); err != nil {
				t.Fatalf("UpdateRivalTag: %s", err)
			}
			if newTag, err := rivalTagService.FindRivalTagByID(tag.ID); err != nil {
				t.Fatalf("FindRivalTagByID: %s", err)
			} else {
				if newTag.TagName != "NEW NAME" {
					t.Fatalf("expected to be %s, got %s", "NEW NAME", newTag.TagName)
				}
				if !newTag.RecordTime.After(currentTime) {
					t.Fatalf("record time is not updated")
				}
				if newTag.Enabled != true {
					t.Fatalf("expected to be true, got %v", newTag.Enabled)
				}
			}
		})
		t.Run("empty value", func(t *testing.T) {
			emptyString := ""
			newEnabled := false
			if err := rivalTagService.UpdateRivalTag(&vo.RivalTagUpdateParam{
				ID:      tag.ID,
				Symbol:  &emptyString,
				Enabled: &newEnabled,
			}); err != nil {
				t.Fatalf("UpdateRivalTag: %s", err)
			}
			if newTag, err := rivalTagService.FindRivalTagByID(tag.ID); err != nil {
				t.Fatalf("FindRivalTagByID: %s", err)
			} else {
				if newTag.TagName == emptyString {
					t.Fatalf("expected to be %s, got %s", tag.TagName, newTag.TagName)
				}
				if newTag.Symbol != emptyString {
					t.Fatalf("expected to be empty string, got %s", newTag.Symbol)
				}
				if newTag.Enabled != newEnabled {
					t.Fatalf("expected to be false, got %v", newTag.Enabled)
				}
			}
		})
	})
}
