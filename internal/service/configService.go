package service

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/charmbracelet/log"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type ConfigService struct {
	db         *gorm.DB
	dirty      chan any
	subscribes []chan<- any
}

func NewConfigService(db *gorm.DB) *ConfigService {
	// NOTE: Some modules in lampghost are working based on config module. And obviously, we don't
	// want to query config current value each time. Therefore, config module implements a notify
	// mechanism.
	ret := &ConfigService{
		db:         db,
		dirty:      make(chan any),
		subscribes: make([]chan<- any, 0),
	}
	go ret.publish()
	return ret
}

// Taking DownloadTaskService, which requires updating the config changes as an example:
// To register the subscription, call Subscribe() function and listen on the channel returned
func (s *ConfigService) Subscribe() <-chan any {
	subscribe := make(chan any)
	s.subscribes = append(s.subscribes, subscribe)
	return subscribe
}

func (s *ConfigService) publish() {
	for {
		notify := <-s.dirty
		for i, ch := range s.subscribes {
			log.Debugf("[ConfigService] sending notification to %d/%v", i, ch)
			ch <- notify
		}
	}
}

func (s *ConfigService) WriteConfig(conf *config.ApplicationConfig) error {
	if err := conf.WriteConfig(); err != nil {
		return err
	}
	s.dirty <- 1
	return nil
}

func (s *ConfigService) QueryLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/Catizard/lampghost/releases/latest")
	if err != nil {
		return "", eris.Wrap(err, "http get github")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", eris.Wrap(err, "read body")
	}
	ret := struct {
		TagName string `json:"tag_name"`
	}{}
	if err := json.Unmarshal(body, &ret); err != nil {
		return "", eris.Wrap(err, "unmarshal")
	}
	return ret.TagName, nil
}

func (s *ConfigService) QueryMetaInfo() (*entity.MetaInfo, error) {
	if version, err := s.QueryLatestVersion(); err != nil {
		return nil, eris.Wrap(err, "query latest version")
	} else {
		return &entity.MetaInfo{
			CurrentVersion: config.VERSION,
			ReleaseVersion: version,
			ClipboardSetup: config.ClipboardSetupFlag,
		}, nil
	}
}
