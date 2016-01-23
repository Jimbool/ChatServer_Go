package configBLL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/dal/configDAL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/config"
	"time"
)

var (
	configObj *config.Config
)

func init() {
	configObj = configDAL.GetConfig()
}

// ====================App相关配置 Begin======================================//
func AppId() string {
	return configObj.AppConfig().AppId
}

func AppName() string {
	return configObj.AppConfig().AppName
}

func AppKey() string {
	return configObj.AppConfig().AppKey
}

// ====================App相关配置 End======================================//

// ====================Socket服务器相关配置 Begin======================================//
func SocketServerAddress() string {
	return configObj.SocketServerConfig().SocketServerAddress()
}

func CheckExpireInterval() time.Duration {
	return configObj.SocketServerConfig().CheckExpireInterval
}

func ClientExpireSeconds() time.Duration {
	return configObj.SocketServerConfig().ClientExpireSeconds
}

func MaxMsgLength() int {
	return configObj.SocketServerConfig().MaxMsgLength
}

// ====================Socket服务器相关配置 End======================================//

// ====================Web服务器相关配置 Begin======================================//

func WebServerAddress() string {
	return configObj.WebServerConfig().WebServerAddress()
}

// ====================Web服务器相关配置 End======================================//
