package difftable

import (
	"encoding/json"
	"io"
	"os"
)

// struct DiffTableData represents one content of a difficult table
// TODO: Sha256 is ignored now, since table may not contain it
type DiffTableData struct {
	Artist    string
	Comment   string
	Level     string
	Lr2BmsId  string `json:"lr2_bmdid"`
	Md5       string
	NameDiff  string
	Title     string
	Url       string `json:"url"`
	UrlDiff   string `json:"url_diff"`
	Sha256    string `json:"-"`
	Lamp      int32  `json:"-"`
	GhostLamp int32  `json:"-"`
}

// Read Difficult Table content from disk.
func ReadDiffTable(filePath string) ([]DiffTableData, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var local []DiffTableData
	err = json.Unmarshal(body, &local)
	if err != nil {
		return nil, err
	}
	return local, nil
}

// Read difficult table content, and return a group split by level
func ReadDiffTableLevelMap(filePath string) (map[string][]DiffTableData, error) {
	rawArr, err := ReadDiffTable(filePath)
	if err != nil {
		return nil, err
	}

	ret := make(map[string][]DiffTableData)
	for _, v := range rawArr {
		if _, ok := ret[v.Level]; !ok {
			ret[v.Level] = make([]DiffTableData, 0)
		}
		ret[v.Level] = append(ret[v.Level], v)
	}
	return ret, nil
}
