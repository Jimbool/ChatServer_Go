/*
项目配置的逻辑处理包，初始化所有的配置内容，其它代码需要配置时都从此包内来获取
包括数据库配置和文件配置
*/
package configBLL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/ChatServer_Go/src/dal/configDAL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/config"
	"time"
)

var (
	appConfigObj          *config.AppConfig
	socketServerConfigObj *config.SocketServerConfig
	webServerConfigObj    *config.WebServerConfig
)

func init() {
	appConfigObj = configDAL.GetAppConfig()
	socketServerConfigObj, webServerConfigObj = configDAL.GetServerConfig(dal.ServerGroupId)
}

// ====================App相关配置 Begin======================================//
func AppId() string {
	return appConfigObj.AppId
}

func AppName() string {
	return appConfigObj.AppName
}

func AppKey() string {
	return appConfigObj.AppKey
}

// ====================App相关配置 End======================================//

// ====================Socket服务器相关配置 Begin======================================//
func SocketServerAddress() string {
	return socketServerConfigObj.SocketServerAddress()
}

func CheckExpireInterval() time.Duration {
	return socketServerConfigObj.CheckExpireInterval
}

func ClientExpireSeconds() time.Duration {
	return socketServerConfigObj.ClientExpireSeconds
}

func MaxMsgLength() int {
	return socketServerConfigObj.MaxMsgLength
}

func MaxHistoryCount() int {
	return socketServerConfigObj.MaxHistoryCount
}

// ====================Socket服务器相关配置 End======================================//

// ====================Web服务器相关配置 Begin======================================//

func WebServerAddress() string {
	return webServerConfigObj.WebServerAddress()
}

// ====================Web服务器相关配置 End======================================//
