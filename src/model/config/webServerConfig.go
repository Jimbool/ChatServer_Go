package config

import (
	"encoding/json"
	"errors"
	"fmt"
)

type WebServerConfig struct {
	ServerHost string // 服务器主机/IP
	ServerPort int    // 服务器端口
}

func NewWebServerConfig(webServerConfigStr string) *WebServerConfig {
	var config *WebServerConfig
	if err := json.Unmarshal([]byte(webServerConfigStr), &config); err != nil {
		panic(errors.New(fmt.Sprintf("反序列化数据出错，原始数据位：%s，错误信息为：%s", webServerConfigStr, err)))
	}

	return config
}

func (config *WebServerConfig) WebServerAddress() string {
	return fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
}
