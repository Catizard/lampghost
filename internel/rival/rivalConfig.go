package rival

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"slices"

	"github.com/Catizard/lampghost/internel/score"
)

const (
	rivalConfigFileName = "rivalConfig.json"
	filePerm            = 0666
)

type RivalTag struct {
	TagName   string
	Generated bool
	TimeStamp int64
}

type RivalInfo struct {
	Name         string
	ScoreLogPath string
	SongDataPath string
	// TODO: ignore tags field temporarily
	Tags     []RivalTag       `json:"-"`
	ScoreLog []score.ScoreLog `json:"-"`
	SongData []score.SongData `json:"-"`
}

func (r *RivalInfo) LoadRivalScoreLog() error {
	scoreLog, err := score.ReadScoreLogFromSqlite(r.ScoreLogPath)
	if err != nil {
		return err
	}
	r.ScoreLog = scoreLog
	return nil
}

func (r *RivalInfo) LoadRivalSongData() error {
	songData, err := score.ReadSongDataFromSqlite(r.SongDataPath)
	if err != nil {
		return err
	}
	r.SongData = songData
	return nil
}

// Save one rival info(or say, meta data) to disk.
func (info *RivalInfo) SaveRivalInfo() error {
	// Sane check before saving rival
	if path.IsAbs(info.ScoreLogPath) {
		return fmt.Errorf("sorry, absolute path is not supported")
	}
	if _, err := os.Stat(info.ScoreLogPath); errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("cannot stat %s on your file system", info.ScoreLogPath)
	}
	if _, err := os.Stat(info.SongDataPath); errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("cannot stat %s on your file system", info.SongDataPath)
	}

	// Read previous data into mermory
	arr, err := ReadRivalInfoFromDisk(rivalConfigFileName)
	if err != nil {
		return err
	}

	// If rival is already added, skip
	if find := slices.ContainsFunc(arr, func(rhs RivalInfo) bool {
		return rhs.Name == info.Name
	}); find {
		return fmt.Errorf("cannot add one rival twice.\nHint: use rival sync to update info instead")
	}

	// If it doesn't exist, append it, then write back
	arr = append(arr, *info)
	// TODO: generate tags here

	newBody, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	os.WriteFile(rivalConfigFileName, newBody, filePerm)
	return nil
}

// Query rival's info by name. Only zero or one result could be match
// Promise that if error is not nil, one rival must be matched
// Warning: If error is not nil, the first result's value has no meaning
func QueryRivalInfo(name string) (RivalInfo, error) {
	// Read disk data into mermory
	arr, err := ReadRivalInfoFromDisk(rivalConfigFileName)
	if err != nil {
		return RivalInfo{}, err
	}

	for _, v := range arr {
		if v.Name == name {
			return v, nil
		}
	}
	return RivalInfo{}, fmt.Errorf("no such a rival named %s", name)
}

// Read rivals data from disk
func ReadRivalInfoFromDisk(path string) ([]RivalInfo, error) {
	// Create file if it doesn't exist
	f, err := os.OpenFile(rivalConfigFileName, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read previous data into mermory
	var prevArray []RivalInfo
	body, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &prevArray)
	if err != nil {
		return nil, err
	}
	return prevArray, nil
}
