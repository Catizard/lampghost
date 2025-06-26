package download

import "fmt"

var _ DownloadSource = (*wriggleDownloadSource)(nil)

var WriggleDownloadSource wriggleDownloadSource = wriggleDownloadSource{
	Meta: DownloadSourceMeta{
		Name:        "wriggle",
		DownloadURL: "https://bms.wrigglebug.xyz/download/package/%s",
	},
}

type wriggleDownloadSource struct {
	Meta DownloadSourceMeta
}

func (d *wriggleDownloadSource) GetMeta() DownloadSourceMeta {
	return d.Meta
}

func (d *wriggleDownloadSource) GetDownloadURLFromMD5(md5 string) (string, string, error) {
	pattern := d.Meta.DownloadURL
	return fmt.Sprintf(pattern, md5), "", nil
}

func (d *wriggleDownloadSource) AllowBatchDownload() bool {
	return false
}
