package service_test

import (
	"testing"

	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/rotisserie/eris"
)

var (
	EmptyScoreLogPath     = "./testdata/empty_scorelog.db"
	EmptySongDataPath     = "./testdata/empty_songdata.db"
	EmptyScoreDataLogPath = "./testdata/empty_scoredatalog.db"

	RealScoreLogPath     = "./testdata/scorelog.db"
	RealSongDataPath     = "./testdata/songdata.db"
	RealScoreDataLogPath = "./testdata/scoredatalog.db"
)

func newNamedEmptyUser(name string, noScoreLog, noSongData, noScoreDataLog bool) *vo.RivalInfoVo {
	scoreLogPath := &EmptyScoreLogPath
	if noScoreLog {
		scoreLogPath = nil
	}
	songDataPath := &EmptySongDataPath
	if noSongData {
		songDataPath = nil
	}
	scoreDataLogPath := &EmptyScoreDataLogPath
	if noScoreDataLog {
		scoreDataLogPath = nil
	}
	return &vo.RivalInfoVo{
		Name:             name,
		ScoreLogPath:     scoreLogPath,
		SongDataPath:     songDataPath,
		ScoreDataLogPath: scoreDataLogPath,
	}
}

func newEmptyUser(noScoreLog, noSongData, noScoreDataLog bool) *vo.RivalInfoVo {
	return newNamedEmptyUser("test", noScoreLog, noSongData, noScoreDataLog)
}

func newNamedRealUser(name string, noScoreLog, noSongData, noScoreDataLog bool) *vo.RivalInfoVo {
	scoreLogPath := &RealScoreLogPath
	if noScoreLog {
		scoreLogPath = nil
	}
	songDataPath := &RealSongDataPath
	if noSongData {
		songDataPath = nil
	}
	scoreDataLogPath := &RealScoreDataLogPath
	if noScoreDataLog {
		scoreDataLogPath = nil
	}
	return &vo.RivalInfoVo{
		Name:             name,
		ScoreLogPath:     scoreLogPath,
		SongDataPath:     songDataPath,
		ScoreDataLogPath: scoreDataLogPath,
	}
}

func newRealUser(noScoreLog, noSongData, noScoreDataLog bool) *vo.RivalInfoVo {
	return newNamedRealUser("test", noScoreLog, noSongData, noScoreDataLog)
}

func skipRealFileTest(noScoreLog, noSongData, noScoreDataLog bool) error {
	if err := database.VerifyLocalDatabaseFilePath(RealScoreLogPath); err != nil && !noScoreLog {
		return eris.New("missing scorelog.db")
	}
	if err := database.VerifyLocalDatabaseFilePath(RealSongDataPath); err != nil && !noSongData {
		return eris.New("missing songdata.db")
	}
	if err := database.VerifyLocalDatabaseFilePath(RealScoreDataLogPath); err != nil && !noScoreDataLog {
		return eris.New("missing scoredatalog.db")
	}
	return nil
}

func TestInitializeMainUser(t *testing.T) {
	t.Run("FastFailOnMissingFilePath", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		rivalInfoService := service.NewRivalInfoService(db)
		missingSongDataPath := newEmptyUser(false, true, true)
		if err := rivalInfoService.InitializeMainUser(missingSongDataPath); err == nil {
			t.Fatalf("should fail on missing songdata file path, but not")
		}
		missingScoreLogPath := newEmptyUser(true, false, false)
		if err := rivalInfoService.InitializeMainUser(missingScoreLogPath); err == nil {
			t.Fatalf("should fail on missing scorelog file path, but not")
		}
	})
	t.Run("EmptyFilesWithoutScoreDataLog", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		rivalInfoService := service.NewRivalInfoService(db)
		emptyFilesWithoutScoreDataLog := newEmptyUser(false, false, true)
		if err := rivalInfoService.InitializeMainUser(emptyFilesWithoutScoreDataLog); err != nil {
			t.Fatalf("failed to initialize main user with empty files(no scoredatalog.db): %s", err)
		}
		mainUser, err := rivalInfoService.QueryMainUser()
		if err != nil {
			t.Fatalf("queryMainUser: %s", err)
		}
		if mainUser.ID != 1 {
			t.Fatalf("assert: mainUser.ID: expected 1, got %d", mainUser.ID)
		}
		if *mainUser.ScoreLogPath != EmptyScoreLogPath {
			t.Fatalf("assert: mainUser.ScoreLogPath: expected %s, got %s", EmptyScoreLogPath, *mainUser.ScoreLogPath)
		}
		if *mainUser.SongDataPath != EmptySongDataPath {
			t.Fatalf("assert: mainUser.SongDataPath: expected %s, got %s", EmptySongDataPath, *mainUser.SongDataPath)
		}
		if mainUser.ScoreDataLogPath != nil {
			t.Fatalf("assert: mainUser.ScoreDataLogPath: expected nil, got %s", *mainUser.ScoreDataLogPath)
		}
	})
	t.Run("EmptyFilesWithScoreDataLog", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		rivalInfoService := service.NewRivalInfoService(db)
		emptyFilesWithoutScoreDataLog := newEmptyUser(false, false, false)
		if err := rivalInfoService.InitializeMainUser(emptyFilesWithoutScoreDataLog); err != nil {
			t.Fatalf("failed to initialize main user with empty files(with scoredatalog.db): %s", err)
		}
		mainUser, err := rivalInfoService.QueryMainUser()
		if err != nil {
			t.Fatalf("queryMainUser: %s", err)
		}
		if mainUser.ID != 1 {
			t.Fatalf("assert: mainUser.ID: expected 1, got %d", mainUser.ID)
		}
		if *mainUser.ScoreLogPath != EmptyScoreLogPath {
			t.Fatalf("assert: mainUser.ScoreLogPath: expected %s, got %s", EmptyScoreLogPath, *mainUser.ScoreLogPath)
		}
		if *mainUser.SongDataPath != EmptySongDataPath {
			t.Fatalf("assert: mainUser.SongDataPath: expected %s, got %s", EmptySongDataPath, *mainUser.SongDataPath)
		}
		if *mainUser.ScoreDataLogPath != EmptyScoreDataLogPath {
			t.Fatalf("assert: mainUser.ScoreDataLogPath: expected %s, got %s", EmptyScoreDataLogPath, *mainUser.ScoreDataLogPath)
		}
	})
	t.Run("RealFiles", func(t *testing.T) {
		var tests = []struct {
			name  string
			input *vo.RivalInfoVo
			skip  error
		}{
			{
				"no scoredatalog.db",
				newRealUser(false, false, true),
				skipRealFileTest(false, false, true),
			},
			{
				"with scoredatalog.db",
				newRealUser(false, false, false),
				skipRealFileTest(false, false, false),
			},
		}
		for _, tt := range tests {
			db, err := database.NewMemoryDatabase()
			if err != nil {
				t.Fatalf("db: %s", err)
			}
			rivalInfoService := service.NewRivalInfoService(db)
			t.Run(tt.name, func(t *testing.T) {
				if tt.skip != nil {
					t.Skipf("skip: %s", tt.skip)
				}
				if err := rivalInfoService.InitializeMainUser(tt.input); err != nil {
					t.Fatalf("initialize main user: %s", err)
				}
			})
		}
	})
}

func TestAddUser(t *testing.T) {
	t.Run("FastFailOnWrongFields", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		rivalInfoService := service.NewRivalInfoService(db)

		var tests = []struct {
			name  string
			input *vo.RivalInfoVo
		}{
			{
				"completely nil",
				nil,
			},
			{
				"name should not be empty",
				newNamedEmptyUser("", false, false, false),
			},
			{
				"scorelog.db path should not be empty",
				newEmptyUser(true, true, false), // no scorelog.db, songdata.db, with scoredatalog.db
			},
			{
				"TODO: seperate songdata.db is not supported",
				newEmptyUser(false, false, false), // with scorelog.db, scoredatalog.db, songdata.db
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := rivalInfoService.AddRivalInfo(tt.input); err == nil {
					t.Fatal("expected error, got nothing")
				}
			})
		}
	})

	t.Run("AddEmptyUser", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		rivalInfoService := service.NewRivalInfoService(db)
		var tests = []struct {
			name  string
			input *vo.RivalInfoVo
		}{
			{
				"without scoredatalog.db",
				newEmptyUser(false, true, true), // with scorelog.db, no scoredatalog.db, no songdata.db
			},
			{
				"with scoredatalog.db",
				newEmptyUser(false, true, false), // with scorelog.db, scoredatalog.db, no songdata.db
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := rivalInfoService.AddRivalInfo(tt.input); err != nil {
					t.Fatalf("add rival: %s", err)
				}
			})

		}
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("FastFailOnWrongFields", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		rivalInfoService := service.NewRivalInfoService(db)
		var tests = []struct {
			name  string
			input *vo.RivalInfoVo
		}{
			{
				"completely nil",
				nil,
			},
			{
				"ID is 0",
				newNamedEmptyUser("", true, true, true),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := rivalInfoService.UpdateRivalInfo(tt.input); err == nil {
					t.Fatalf("expected error, got nothing")
				}
			})
		}
	})

	t.Run("UpdatePlainFields", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		rivalInfoService := service.NewRivalInfoService(db)
		emptyMainUser := newEmptyUser(false, false, false)
		if err := rivalInfoService.InitializeMainUser(emptyMainUser); err != nil {
			t.Fatalf("initialize main user: %s", err)
		}
		updateParam := newNamedEmptyUser("NEWNAME", true, true, true)
		updateParam.SongDataPath = emptyMainUser.SongDataPath
		updateParam.ID = 1
		if err := rivalInfoService.UpdateRivalInfo(updateParam); err != nil {
			t.Fatalf("update rival: %s", err)
		}
		currState, err := rivalInfoService.QueryMainUser()
		if err != nil {
			t.Fatalf("query: %s", err)
		}
		if currState == nil {
			t.Fatalf("query: no user found")
		}
		if currState.Name != "NEWNAME" {
			t.Error("update rival: name doesn't be modified")
		}
		if !currState.UpdatedAt.After(currState.CreatedAt) {
			t.Error("update rival: last update time wasn't being updated")
		}
	})

	// Some fields of the rival should never be updated directly through UpdateRivalInfo interface
	t.Run("NoDirectlyWriteFields", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		rivalInfoService := service.NewRivalInfoService(db)
		emptyUser := newEmptyUser(false, true, false)
		if err := rivalInfoService.AddRivalInfo(emptyUser); err != nil {
			t.Fatalf("initialize main user: %s", err)
		}
		updateParam := newNamedEmptyUser("NEWNAME", true, true, true)
		updateParam.ID = 1
		// These fields should never be updated directly
		updateParam.PlayCount = 1
		updateParam.MainUser = true
		if err := rivalInfoService.UpdateRivalInfo(updateParam); err != nil {
			t.Fatalf("update rival: %s", err)
		}
		currState, err := rivalInfoService.FindRivalInfoByID(1)
		if err != nil {
			t.Fatalf("query: %s", err)
		}
		if currState == nil {
			t.Fatalf("query: no user found")
		}
		if currState.PlayCount != 0 {
			t.Error("directly writes PlayCount detected")
		}
		if currState.MainUser != false {
			t.Errorf("directly writes MainUser detected")
		}
	})
}
