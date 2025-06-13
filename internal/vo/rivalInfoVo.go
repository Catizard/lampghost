package vo

import (
	"path"

	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type RivalInfoVo struct {
	gorm.Model
	Name             string
	ScoreLogPath     *string
	SongDataPath     *string
	ScoreDataLogPath *string
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
		ScoreLogPath:     rivalInfo.ScoreLogPath,
		SongDataPath:     rivalInfo.SongDataPath,
		ScoreDataLogPath: rivalInfo.ScoreDataLogPath,
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
	ImportStrategy         string // constants: "directory" | "separate"
	BeatorajaDirectoryPath string
	PlayerDirectory        string
	ScoreLogPath           *string
	SongDataPath           *string
	ScoreDataLogPath       *string
}

func (rivalInfo *InitializeRivalInfoVo) Into() *entity.RivalInfo {
	if rivalInfo.ImportStrategy == "separate" {
		return &entity.RivalInfo{
			Name:             rivalInfo.Name,
			SongDataPath:     rivalInfo.SongDataPath,
			ScoreLogPath:     rivalInfo.ScoreLogPath,
			ScoreDataLogPath: rivalInfo.ScoreDataLogPath,
		}
	} else {
		// NOTE: we don't need to do much verification here
		log.Debugf("beatoraja: %v", rivalInfo.BeatorajaDirectoryPath)
		songdataPath := path.Join(rivalInfo.BeatorajaDirectoryPath, "songdata.db")
		log.Debugf("songdata: %v", songdataPath)
		saveFilePath := path.Join(rivalInfo.BeatorajaDirectoryPath, "player", rivalInfo.PlayerDirectory)
		scorelogPath := path.Join(saveFilePath, "scorelog.db")
		scoredatalogPath := path.Join(saveFilePath, "scoredatalog.db")
		return &entity.RivalInfo{
			Name:             rivalInfo.Name,
			SongDataPath:     &songdataPath,
			ScoreLogPath:     &scorelogPath,
			ScoreDataLogPath: &scoredatalogPath,
		}
	}
}
