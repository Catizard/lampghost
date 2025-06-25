package config

type DownloadSourceMeta struct {
	Name        string
	DownloadURL string
}

type DownloadSource interface {
	GetDownloadURLFromMD5(md5 string) string
	GetMeta() DownloadSourceMeta
	AllowBatchDownload() bool
}
