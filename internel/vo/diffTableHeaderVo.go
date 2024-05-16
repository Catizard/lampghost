package vo

type DiffTableHeader struct {
	DataUrl     string `json:"data_url"`
	LastUpdate  string `json:"last_update"`
	Name        string
	OriginalUrl string `json:"original_url"`
	Symbol      string
}
