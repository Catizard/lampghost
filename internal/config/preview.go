package config

import "fmt"

var DefaultPreviewSource PreviewSource = &SayakaPreviewSource
var previewSources map[string]PreviewSource = map[string]PreviewSource{
	"sayaka": &SayakaPreviewSource,
	"konmai": &KonmaiPreviewSource,
}

type PreviewSource interface {
	GetName() string
	GetPreviewURLByMd5(string) (string, error)
}

type sayakaPreviewSource struct{}

func (s *sayakaPreviewSource) GetName() string { return "sayaka" }
func (s *sayakaPreviewSource) GetPreviewURLByMd5(md5 string) (string, error) {
	return fmt.Sprintf("https://bms-score-viewer.pages.dev/view?md5=%s", md5), nil
}

var SayakaPreviewSource sayakaPreviewSource = sayakaPreviewSource{}

type konmaiPreviewSource struct{}

func (s *konmaiPreviewSource) GetName() string { return "konmai" }
func (s *konmaiPreviewSource) GetPreviewURLByMd5(md5 string) (string, error) {
	return fmt.Sprintf("http://bms.alvorna.com/bms/score/?md5=%s", md5), nil
}

var KonmaiPreviewSource konmaiPreviewSource = konmaiPreviewSource{}

func GetPreviewSource(name string) PreviewSource {
	return previewSources[name]
}
