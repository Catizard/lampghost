package difftable

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/Catizard/lampghost/internel/remote"
)

const (
	headerFileName = "tableHeader.json"
)

type DiffTableHeader struct {
	DataUrl     string `json:"data_url"`
	LastUpdate  string `json:"last_update"`
	Name        string
	OriginalUrl string `json:"original_url"`
	Symbol      string `json:"Symbol"`
	Alias       string
}

// Add a difficult table header(or say, meta data) and related song data to disk.
// Before calling this function, data.json should't exist on disk, otherwise beaviour is undefined.
// All difficult table headers' info would be saved in one single file.
func (header *DiffTableHeader) AddDiffTable() error {
	// 1. Sync table header file
	headers, err := QueryDifficultTableHeaderStrictly(header.Name)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		return fmt.Errorf("%s difficult table has already been added.\nHint: use table sync %s to sync table's data\nIf you believe this is an error, type 'table del %s' to remove", header.Name, header.Name, header.Name)
	}

	// Append new one, then write everything back
	headers = append(headers, *header)

	newBody, err := json.Marshal(headers)
	if err != nil {
		panic(err)
	}
	os.WriteFile(headerFileName, newBody, fs.FileMode(os.O_WRONLY))

	// 2. Create data file
	fileName := header.getDataJsonFileName()
	// If data.json is already here, do nothing
	if _, err := os.Stat(fileName); err == nil {
		panic(fmt.Errorf("%s has been already exists, if you want to update data.json, use sync command instead", fileName))
	} else if errors.Is(err, fs.ErrExist) {
		// unexpected...
		panic(err)
	}
	return saveTableData(header.getDataJsonFileName(), header.DataUrl)
}

// Query table headers on disk.
// Returned when name or alias name is matched
func QueryDifficultTableHeaderLoosely(name string) ([]DiffTableHeader, error) {
	return QueryDifficultTableHeader(func(dth DiffTableHeader) bool {
		return dth.Name == name || dth.Alias == name
	})
}

// Query table headers on disk.
// Returned only when name is matched
func QueryDifficultTableHeaderStrictly(name string) ([]DiffTableHeader, error) {
	return QueryDifficultTableHeader(func(dth DiffTableHeader) bool {
		return dth.Name == name
	})
}

// Prototype of difficult table header's query function.
// Query headers from disk, return array of headers when succeed
func QueryDifficultTableHeader(equalf func(DiffTableHeader) bool) ([]DiffTableHeader, error) {
	// Read all data from disk.
	var local []DiffTableHeader
	remote.FetchJsonFromFile(headerFileName, &local)

	res := make([]DiffTableHeader, 0)
	for _, v := range local {
		if equalf(v) {
			res = append(res, v)
		}
	}
	return res, nil
}

// Query but panic when multiple result matched
// None matched is also take as an error
// Simple wrap of QueryDifficultTableHeader
func QueryDifficultTableHeaderExactlyOne(name string) DiffTableHeader {
	rawArr, err := QueryDifficultTableHeaderLoosely(name)
	if err != nil {
		panic(err)
	}
	if len(rawArr) == 0 {
		panic("no difficult table header found")
	}
	if len(rawArr) > 1 {
		panic("multiple difficult table matched")
	}
	return rawArr[0]
}

func (header *DiffTableHeader) SyncDifficultTable() error {
	// TODO: header's info might changed, sync header part also

	// sync data.json, propagate any error
	return saveTableData(header.getDataJsonFileName(), header.DataUrl)
}

// Download specified difficult table's data to disk
// File would be overwrite if it already exist
func saveTableData(fileName string, dataUrl string) error {
	// 1. Create data file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// 2. Download to file
	// TODO: if dataUrl is not start with http...
	resp, err := http.Get(dataUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(file, resp.Body)
	return nil
}

// Return difficult table header's data json file name
func (header *DiffTableHeader) getDataJsonFileName() string {
	return header.Name + ".json"
}
