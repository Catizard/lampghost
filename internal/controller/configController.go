package controller

import (
	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/result"
	"github.com/Catizard/lampghost_wails/internal/service"
)

type ConfigController struct {
	service *service.ConfigService
}

func NewConfigController(configService *service.ConfigService) *ConfigController {
	return &ConfigController{
		service: configService,
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
	if err := ctl.service.WriteConfig(conf); err != nil {
		return result.NewErrorMessage(err)
	}
	return result.SUCCESS
}
