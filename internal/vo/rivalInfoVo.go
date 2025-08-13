package vo

import (
	"path"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type RivalInfoVo struct {
	gorm.Model
	Name             string
	Type             string
	ScoreLogPath     *string
	SongDataPath     *string
	ScoreDataLogPath *string
	ScoreDataPath    *string
	PlayCount        int
	MainUser         bool
	// Below fields cannot be updated by `UpdateRivalInfo`
	LockTagID     uint
	ReverseImport int

	Pagination     *entity.Page
	Locale         *string // only passed at initialized phase
	IgnoreMainUser bool    // equivalent to ID != 1
}

func (rivalInfo *RivalInfoVo) Entity() *entity.RivalInfo {
	return &entity.RivalInfo{
		Model: gorm.Model{
			ID:        rivalInfo.ID,
			CreatedAt: rivalInfo.CreatedAt,
			UpdatedAt: rivalInfo.UpdatedAt,
			DeletedAt: rivalInfo.DeletedAt,
		},
		Name:             rivalInfo.Name,
		Type:             rivalInfo.Type,
		ScoreLogPath:     rivalInfo.ScoreLogPath,
		SongDataPath:     rivalInfo.SongDataPath,
		ScoreDataLogPath: rivalInfo.ScoreDataLogPath,
		ScoreDataPath:    rivalInfo.ScoreDataPath,
		PlayCount:        rivalInfo.PlayCount,
		MainUser:         rivalInfo.MainUser,
		LockTagID:        rivalInfo.LockTagID,
		ReverseImport:    rivalInfo.ReverseImport,
	}
}

// For initialization usage, could be converted
// into a RivalInfo.
type InitializeRivalInfoVo struct {
	Name                   string
	Locale                 *string
	ImportStrategy         string `json:"ImportStrategy"` // constants: "directory" | "separate" | "LR2"
	BeatorajaDirectoryPath string
	PlayerDirectory        string
	ScoreLogPath           *string
	SongDataPath           *string
	ScoreDataLogPath       *string
	ScoreDataPath          *string
	BMSDirectories         []string `json:"BMSDirectories"` // replacement of SongDataPath
}

func (rivalInfo *InitializeRivalInfoVo) Into() *entity.RivalInfo {
	switch rivalInfo.ImportStrategy {
	case "separate":
		return &entity.RivalInfo{
			Name:             rivalInfo.Name,
			Type:             entity.RIVAL_TYPE_BEATORAJA,
			SongDataPath:     rivalInfo.SongDataPath,
			ScoreLogPath:     rivalInfo.ScoreLogPath,
			ScoreDataLogPath: rivalInfo.ScoreDataLogPath,
			ScoreDataPath:    rivalInfo.ScoreDataPath,
		}
	case "directory":
		// NOTE: we don't need to do much verification here
		log.Debugf("beatoraja: %v", rivalInfo.BeatorajaDirectoryPath)
		songdataPath := path.Join(rivalInfo.BeatorajaDirectoryPath, "songdata.db")
		log.Debugf("songdata: %v", songdataPath)
		saveFilePath := path.Join(rivalInfo.BeatorajaDirectoryPath, "player", rivalInfo.PlayerDirectory)
		scorelogPath := path.Join(saveFilePath, "scorelog.db")
		scoredatalogPath := path.Join(saveFilePath, "scoredatalog.db")
		scoredataPath := path.Join(saveFilePath, "score.db")
		return &entity.RivalInfo{
			Name:             rivalInfo.Name,
			Type:             entity.RIVAL_TYPE_BEATORAJA,
			SongDataPath:     &songdataPath,
			ScoreLogPath:     &scorelogPath,
			ScoreDataLogPath: &scoredatalogPath,
			ScoreDataPath:    &scoredataPath,
		}
	case "LR2":
		selfGeneratedSongDataPath := config.GetSelfGeneratedSongDataPath()
		return &entity.RivalInfo{
			Name:             rivalInfo.Name,
			Type:             entity.RIVAL_TYPE_LR2,
			ScoreLogPath:     rivalInfo.ScoreLogPath,
			ScoreDataLogPath: rivalInfo.ScoreDataLogPath,
			SongDataPath:     &selfGeneratedSongDataPath,
			ScoreDataPath:    rivalInfo.ScoreLogPath,
		}
	default:
		log.Errorf("unexpected import strategy: %s", rivalInfo.ImportStrategy)
		return nil
	}
}

// Indicates if we need to reload specific file
type RivalFileReloadInfoVo struct {
	SongData     bool
	ScoreLog     bool
	ScoreDataLog bool
	ScoreData    bool
}
