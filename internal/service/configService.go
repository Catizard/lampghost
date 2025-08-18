package service

import (
	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/charmbracelet/log"
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
