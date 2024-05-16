package remote

import (
	"encoding/json"
	"io"
	"net/http"
)

func FetchJson(url string, v interface{}) {
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
