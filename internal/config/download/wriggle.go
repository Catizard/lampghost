package download

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rotisserie/eris"
)

var _ DownloadSource = (*wriggleDownloadSource)(nil)

var WriggleDownloadSource wriggleDownloadSource = wriggleDownloadSource{
	Meta: DownloadSourceMeta{
		Name:         "wriggle",
		DownloadURL:  "https://bms.wrigglebug.xyz/download/package/%s",
		MetaQueryURL: "https://bms.wrigglebug.xyz/api/package/%s",
	},
}

type wriggleDownloadSource struct {
	Meta DownloadSourceMeta
}

func (d *wriggleDownloadSource) GetMeta() DownloadSourceMeta {
	return d.Meta
}

func (d *wriggleDownloadSource) GetDownloadURLFromMD5(md5 string) (downloadInfo DownloadInfo, err error) {
	metaQueryURL := fmt.Sprintf(d.Meta.MetaQueryURL, md5)
	resp, err := http.Get(metaQueryURL)
	if err != nil {
		err = eris.Wrap(err, "get meta")
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		err = eris.Wrap(err, "http read")
		return
	}
	if strings.HasPrefix(string(b), "404") {
		err = eris.Errorf("404 NOT FOUND")
		return
	}
	var result mResp
	if err = json.Unmarshal(b, &result); err != nil {
		err = eris.Wrapf(err, "failed to unmarshal result: %s", string(b))
		return
	}
	if result.Error != "" {
		err = eris.New(result.Error)
		return
	}
	return DownloadInfo{
		DownloadURL:  fmt.Sprintf(d.Meta.DownloadURL, md5),
		UniqueSymbol: result.FileHash,
		FileName:     result.Name,
	}, nil
}

func (d *wriggleDownloadSource) AllowBatchDownload() bool {
	return true
}

// Wriggle server models
type mResp struct {
	DiskSizeBytes uint32 `json:"disk_size_bytes"`
	FileHash      string `json:"file_hash"`
	Hashes        map[string]string
	ModTime       string `json:"mod_time"`
	Name          string
	Error         string
}
