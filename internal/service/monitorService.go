package service

import (
	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
)

// For monitoring the scorelog.db file changes
// NOTE: Currently, MonitorService doesn't listen the changes of config file
type MonitorService struct {
	watcher               *fsnotify.Watcher
	currentWatching       *string
	notifySyncChan        chan<- any
	subscribeConfigChange <-chan any
	enableAutoReload      bool
}

func NewMonitorService(config *config.ApplicationConfig) (*MonitorService, <-chan any) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	notifySyncChan := make(chan any)

	ret := &MonitorService{
		watcher:        watcher,
		notifySyncChan: notifySyncChan,
	}

	if config.EnableAutoReload != 0 {
		ret.enableAutoReload = true
	} else {
		ret.enableAutoReload = false
	}

	go ret.listen()

	return ret, notifySyncChan
}

func (s *MonitorService) listen() {
	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				log.Errorf("[MonitorService] unexpected error, exited")
				return
			}
			log.Debugf("[MonitorService] event: %v", event)
			if s.enableAutoReload && event.Has(fsnotify.Write) {
				log.Debugf("[MonitorService] sending notification to RivalInfoService")
				s.notifySyncChan <- new(any)
			}
		case err, ok := <-s.watcher.Errors:
			if !ok {
				log.Errorf("[MonitorService] unexpected error, exited")
				return
			}
			log.Errorf("error: %s", err)
		case <-s.subscribeConfigChange:
			conf, err := config.ReadConfig()
			if err != nil {
				log.Errorf("cannot read config: %s", err)
			} else {
				if conf.EnableAutoReload != 0 {
					s.enableAutoReload = true
				} else {
					s.enableAutoReload = false
				}
			}
		}
	}
}

// Set scorelog.db file path to listen
func (s *MonitorService) SetScoreLogFilePath(fp string) error {
	log.Debugf("[MonitorService] trying to set scorelog.db file path to %s", fp)
	if err := database.VerifyLocalDatabaseFilePath(fp); err != nil {
		return err
	}
	if s.currentWatching != nil {
		s.watcher.Remove(*s.currentWatching)
	}
	s.currentWatching = &fp
	return s.watcher.Add(fp)
}
