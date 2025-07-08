package download

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var KonmaiDownloadSource konmaiDownloadSource = konmaiDownloadSource{
	Meta: DownloadSourceMeta{
		Name:        "konmai",
		DownloadURL: "https://bms.alvorna.com/api/hash?md5=%s",
	},
}

type konmaiDownloadSource struct {
	Meta DownloadSourceMeta
}

func (d *konmaiDownloadSource) GetMeta() DownloadSourceMeta {
	return d.Meta
}

func (d *konmaiDownloadSource) GetDownloadURLFromMD5(md5 string) (downloadInfo DownloadInfo, err error) {
	metaQueryURL := fmt.Sprintf(d.Meta.DownloadURL, md5)
	resp, err := http.Get(metaQueryURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var result kResp
	if err = json.Unmarshal(b, &result); err != nil {
		return
	}
	if result.Result != "success" {
		err = fmt.Errorf("konmai: %s", result.Msg)
		return
	}
	if result.Data.SongURL == "" {
		err = fmt.Errorf("Not Found")
		return
	}
	return DownloadInfo{
		DownloadURL:  result.Data.SongURL,
		UniqueSymbol: result.Data.SongURL,
		FileName:     fmt.Sprintf("%s.7z", result.Data.SongName),
	}, nil
}

func (d *konmaiDownloadSource) AllowBatchDownload() bool {
	return true
}

// Konmai server models
type kResp struct {
	Result string
	Msg    string
	Chart  string
	Data   kChart
}

type kChart struct {
	ChartName string `json:"chart_name"`
	Md5       string
	Sha256    string
	SongName  string `json:"song_name"`
	SongURL   string `json:"song_url"`
}
