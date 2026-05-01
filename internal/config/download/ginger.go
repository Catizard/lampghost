package download

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rotisserie/eris"
)

var _ DownloadSource = (*gingerDownloadSource)(nil)

var GingerDownloadSource gingerDownloadSource = gingerDownloadSource{
	Meta: DownloadSourceMeta{
		Name:         "ginger",
		MetaQueryURL: "https://gingerrush.com/api/v1/files/package/%s",
	},
}

type gingerDownloadSource struct {
	Meta DownloadSourceMeta
}

func (d *gingerDownloadSource) GetMeta() DownloadSourceMeta {
	return d.Meta
}

func (d *gingerDownloadSource) GetDownloadURLFromMD5(md5 string) (downloadInfo DownloadInfo, err error) {
	metaQueryURL := fmt.Sprintf(d.Meta.MetaQueryURL, md5)
	resp, err := d.queryPackage(metaQueryURL)
	if err != nil {
		return DownloadInfo{}, err
	}
	return DownloadInfo{
		DownloadURL:  resp.DownloadURL,
		UniqueSymbol: resp.DownloadURL,
		FileName:     resp.FileName,
	}, nil
}

func (d *gingerDownloadSource) queryPackage(metaQueryURL string) (*gResp, error) {
	resp, err := http.Get(metaQueryURL)
	if err != nil {
		return nil, eris.Wrap(err, "get meta")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, eris.Errorf("error code: %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, eris.Wrap(err, "http read")
	}
	if strings.HasPrefix(string(b), "404") {
		return nil, eris.Errorf("404 NOT FOUND")
	}
	var result gResp
	if err = json.Unmarshal(b, &result); err != nil {
		return nil, eris.Wrapf(err, "failed to unmarshal result: %s", string(b))
	}
	return &result, nil
}

func (d *gingerDownloadSource) AllowBatchDownload() bool {
	return true
}

// Ginger server models
type gResp struct {
	ShardMD5    string `json:"shardMD5"`
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	DirectoryID string `json:"directoryID"`
	MD5s        string `json:"md5s"`
	DownloadURL string `json:"downloadURL"`
}
