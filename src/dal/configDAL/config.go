/*
系统配置的逻辑处理包
*/
package configDAL

import (
	"errors"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/ChatServer_Go/src/model/config"
)

func GetConfig() *config.Config {
	sql := "SELECT AppId, AppName, AppKey, SocketServerConfig, WebServerConfig FROM config;"

	rows, err := dal.DB().Query(sql)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var appId string
		var appName string
		var appKey string
		var socketServerConfig string
		var webServerConfig string
		err := rows.Scan(&appId, &appName, &appKey, &socketServerConfig, &webServerConfig)
		if err != nil {
			panic(err)
		}

		return config.NewConfig(appId, appName, appKey, socketServerConfig, webServerConfig)
	}

	panic(errors.New("未找到config配置内容"))
}
