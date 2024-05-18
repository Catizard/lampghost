package difftable

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
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
	// create table header file if necessary
	// TODO: A constant variable file?
	var file *os.File

	var headers []DiffTableHeader
	// race?
	if _, err := os.Stat(headerFileName); errors.Is(err, fs.ErrNotExist) {
		_, err = os.Create(headerFileName)
		if err != nil {
			panic(err)
		}
	} else {
		file, err = os.Open(headerFileName)
		if err != nil {
			panic(err)
		}
		body, err := io.ReadAll(file)
		log.Printf("body=%s", body)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(body, &headers)
		if err != nil {
			panic(err)
		}
		log.Printf("len(headers)=%v", len(headers))
	}

	// if new one already exist
	for _, exHeader := range headers {
		if exHeader.Name == header.Name {
			// TODO: del command?
			return fmt.Errorf("%s difficult table has already been added.\nHint: use table sync %s to sync table's data\nIf you believe this is an error, type 'table del %s' to remove", header.Name, header.Name, header.Name)
		}
	}

	// append new one, then write everything back
	headers = append(headers, *header)

	newBody, err := json.Marshal(headers)
	if err != nil {
		panic(err)
	}
	os.WriteFile(headerFileName, newBody, fs.FileMode(os.O_WRONLY))

	// 2. Create data file
	fileName := fmt.Sprintf("%s.json", header.Name)
	// If data.json is already here, do nothing
	if _, err := os.Stat(fileName); err == nil {
		panic(fmt.Errorf("%s is already exists, if you want to update data.json, use sync command instead", fileName))
	} else if errors.Is(err, fs.ErrExist) {
		// unexpected...
		panic(err)
	}
	return saveTableData(header.getDataJsonFileName(), header.DataUrl)
}

// Query loaded table headers on disk.
// Use name to identify, match on alias or original name.
// Matched result count could be more than 1.
func QueryDifficultTableHeader(name string) ([]DiffTableHeader, error) {
	// Read all data from disk.
	f, err := os.Open(headerFileName)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var local []DiffTableHeader
	err = json.Unmarshal(body, &local)
	if err != nil {
		return nil, err
	}

	res := make([]DiffTableHeader, 0)

	for _, v := range local {
		if v.Name == name || v.Alias == name {
			res = append(res, v)
		}
	}
	return res, nil
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
