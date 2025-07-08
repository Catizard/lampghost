package download

var DefaultDownloadSource DownloadSource = &KonmaiDownloadSource
var downloadSources map[string]DownloadSource = map[string]DownloadSource{
	"wriggle": &WriggleDownloadSource,
	"konmai":  &KonmaiDownloadSource,
}

type DownloadSourceMeta struct {
	Name         string
	DownloadURL  string
	MetaQueryURL string
}

type DownloadSource interface {
	GetDownloadURLFromMD5(string) (DownloadInfo, error)
	GetMeta() DownloadSourceMeta
	AllowBatchDownload() bool
}

type DownloadInfo struct {
	DownloadURL  string
	UniqueSymbol string
	FileName     string
}

func GetDownloadSource(name string) DownloadSource {
	return downloadSources[name]
}
