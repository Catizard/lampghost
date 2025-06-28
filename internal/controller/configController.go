package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Catizard/lampghost_wails/internal/config"
	"github.com/Catizard/lampghost_wails/internal/config/download"
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

func (ctl *ConfigController) QueryLatestVersion() result.RtnMessage {
	resp, err := http.Get("https://api.github.com/repos/Catizard/lampghost/releases/latest")
	if err != nil {
		return result.NewErrorMessage(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result.NewErrorMessage(err)
	}
	ret := struct {
		TagName string `json:"tag_name"`
	}{}
	if err := json.Unmarshal(body, &ret); err != nil {
		return result.NewErrorMessage(err)
	}
	if ret.TagName != config.VERSION {
		return result.NewRtnMessage(result.SUCCESS.Code, fmt.Sprintf("Release: %s; Using %s", ret.TagName, config.VERSION))
	}
	return result.NewRtnMessage(result.SUCCESS.Code, "Using the latest version")
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

func (ctl *ConfigController) QueryCurrentDownloadSource() result.RtnData {
	if conf, err := config.ReadConfig(); err != nil {
		return result.NewErrorData(err)
	} else {
		downloadSource := download.GetDownloadSource(conf.DownloadSite)
		return result.NewRtnData(downloadSource)
	}
}

func (ctl *ConfigController) GENERATOR_DOWNLOAD_SOURCE() *download.DownloadSource { return nil }
