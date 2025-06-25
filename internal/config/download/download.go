package download

var DefaultDownloadSource DownloadSource = &KonmaiDownloadSource
var downloadSources map[string]DownloadSource = map[string]DownloadSource{
	"wriggle": &WriggleDownloadSource,
	"konmai":  &KonmaiDownloadSource,
}

type DownloadSourceMeta struct {
	Name        string
	DownloadURL string
}

type DownloadSource interface {
	GetDownloadURLFromMD5(string) (string, error)
	GetMeta() DownloadSourceMeta
	AllowBatchDownload() bool
}

func GetDownloadSource(name string) DownloadSource {
	return downloadSources[name]
}
