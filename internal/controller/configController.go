package controller

import (
	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
)

type ConfigController struct {
	rivalInfoService *service.RivalInfoService
}

func NewConfigController(rivalInfoService *service.RivalInfoService) *ConfigController {
	return &ConfigController{
		rivalInfoService: rivalInfoService,
	}
}

func (ctl *ConfigController) ReadConfig() result.RtnData {
	if conf, err := config.ReadConfig(); err != nil {
		return result.NewErrorData(err)
	} else {
		return result.NewRtnData(conf)
	}
}

func (ctl *ConfigController) WriteConfig(conf *config.ApplicationConfig) result.RtnMessage {
	prevConf, err := config.ReadConfig()
	if err != nil {
		return result.NewErrorMessage(err)
	}
	if shouldReload(prevConf, conf) {
		mainUser, err := ctl.rivalInfoService.QueryMainUser()
		if err != nil {
			return result.NewErrorMessage(err)
		}
		mainUser.ScoreLogPath = &conf.ScorelogFilePath
		mainUser.SongDataPath = &conf.SongdataFilePath
		// TODO: mainUser.scorePath = &conf.ScoreFilePath
		if err := ctl.rivalInfoService.SyncRivalScoreLog(mainUser); err != nil {
			return result.NewErrorMessage(err)
		}
	}
	// TODO: what if the player is reloaded successfully but the file isn't?
	if err := conf.WriteConfig(); err != nil {
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}

func shouldReload(prevConf *config.ApplicationConfig, newConf *config.ApplicationConfig) bool {
	if prevConf.ScorelogFilePath != newConf.ScorelogFilePath {
		return true
	}
	if prevConf.SongdataFilePath != newConf.SongdataFilePath {
		return true
	}
	if prevConf.ScoreFilePath != newConf.ScoreFilePath {
		return true
	}
	return false
}
