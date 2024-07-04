package common

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

// Simple wrapper for reading json from url
// Panic when error occured
func FetchJson(url string, v interface{}) {
	log.Infof("Fetching json from %s", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		panic(err)
	}
}

// Simple wrapper for reading json from file.
// If file doesn't exist, create it.
// Panic when error occured.
func FetchJsonFromFile(filePath string, v interface{}) {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	body, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		panic(err)
	}
}
