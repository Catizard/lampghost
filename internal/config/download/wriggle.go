package download

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var result mResp
	if err = json.Unmarshal(b, &result); err != nil {
		return
	}
	if result.Error != "" {
		err = eris.Errorf(result.Error)
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
	DiskSizeBytes uint32
	FileHash      string
	Hashes        map[string]string
	ModTime       string
	Name          string
	Error         string
}
