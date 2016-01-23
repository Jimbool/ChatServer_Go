package config

import (
	"encoding/json"
)

type Config struct {
	appId                 string              //App唯一标识
	appName               string              //App名称
	appKey                string              //App加密Key
	appConfig             *AppConfig          //应用配置对象
	socketServerConfigStr string              //Socket服务器配置字符串
	socketServerConfig    *SocketServerConfig //Socket服务器配置对象
	webServerConfigStr    string              //Web服务器配置字符串
	webServerConfig       *WebServerConfig    //Web服务器配置对象
}

func NewConfig(appId, appName, appKey, socketServerConfig, webServerConfig string) *Config {
	config := &Config{
		appId:                 appId,
		appName:               appName,
		appKey:                appKey,
		socketServerConfigStr: socketServerConfig,
		webServerConfigStr:    webServerConfig,
	}

	config.initAppConfig()
	config.initSocketServerConfig()
	config.initWebServerConfig()

	return config
}

func (config *Config) initAppConfig() {
	config.appConfig = NewAppConfig(config.appId, config.appName, config.appKey)
}

func (config *Config) initSocketServerConfig() {
	if err := json.Unmarshal([]byte(config.socketServerConfigStr), &config.socketServerConfig); err != nil {
		panic(err)
	}
}

func (config *Config) initWebServerConfig() {
	if err := json.Unmarshal([]byte(config.webServerConfigStr), &config.webServerConfig); err != nil {
		panic(err)
	}
}

func (config *Config) AppConfig() *AppConfig {
	return config.appConfig
}

func (config *Config) SocketServerConfig() *SocketServerConfig {
	return config.socketServerConfig
}

func (config *Config) WebServerConfig() *WebServerConfig {
	return config.webServerConfig
}
