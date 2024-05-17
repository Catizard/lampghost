package vo

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

type DiffTableHeader struct {
	DataUrl     string `json:"data_url"`
	LastUpdate  string `json:"last_update"`
	Name        string
	OriginalUrl string `json:"original_url"`
	Symbol      string `json:"Symbol"`
	Alias       string
}

// Add a difficult table header(or say, meta data) and related song data to disk.
// Before calling this function, data.json should't exist on disk, beaviour is undefined.
// All difficult table headers would locate in one single file.If the file hasn't exist, it would be created.
func (header *DiffTableHeader) AddDiffTable() error {
	// 1. Sync table header file
	// create table header file if necessary
	// TODO: A constant variable file?
	headerFileName := "tableHeader.json"
	var file *os.File

	var headers []DiffTableHeader
	// race?
	if _, err := os.Stat(headerFileName); errors.Is(err, fs.ErrNotExist) {
		file, err = os.Create(headerFileName)
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
	// Or, add this difficult table
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	// 3. Download to file
	// TODO: if dataUrl is not start with http...
	resp, err := http.Get(header.DataUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(file, resp.Body)
	return nil
}
