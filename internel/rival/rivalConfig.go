package rival

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
)

const (
	rivalConfigFileName = "rivalConfig.json"
)

type RivalInfo struct {
	Name string
	ScoreLogPath string
}

// Add one rival info(or say, meta data) to disk.
func AddRivalInfo(info *RivalInfo) error {
	// Before we do anything, check if scorelog filepath exist
	if path.IsAbs(info.ScoreLogPath) {
		panic("Sorry, absolute path is not supported.")
	}
	if _, err := os.Stat(info.ScoreLogPath); errors.Is(err, fs.ErrNotExist) {
		panic(fmt.Errorf("cannot find %s on your file system", info.ScoreLogPath))
	}
	var prevArray []RivalInfo

	if _, err := os.Stat(rivalConfigFileName); errors.Is(err, fs.ErrNotExist) {
		_, err := os.Create(rivalConfigFileName)
		if err != nil {
			return err
		}
	} else {
		file, err := os.Open(rivalConfigFileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		body, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, &prevArray)
		if err != nil {
			return err
		}
	}

	find := false
	for i, v := range prevArray {
		if v.Name == info.Name {
			// We don't have to remove it, we can override it
			// But this kind of sucks
			prevArray[i].Name = info.Name
			prevArray[i].ScoreLogPath = info.ScoreLogPath
			find = true
		}
	}
	
	// If it doesn't exist, append it
	if !find {
		prevArray = append(prevArray, *info)
	}

	newBody, err := json.Marshal(prevArray)
	if err != nil {
		return err
	}
	os.WriteFile(rivalConfigFileName, newBody, fs.FileMode(os.O_WRONLY))
	return nil
}