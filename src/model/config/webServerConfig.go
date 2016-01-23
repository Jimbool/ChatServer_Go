package config

import (
	"fmt"
)

type WebServerConfig struct {
	ServerHost string // 服务器主机/IP
	ServerPort int    // 服务器端口
}

func (config *WebServerConfig) WebServerAddress() string {
	return fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
}
