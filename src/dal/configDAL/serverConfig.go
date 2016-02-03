package configDAL

import (
	"errors"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/ChatServer_Go/src/model/config"
)

func GetServerConfig(id int) (*config.SocketServerConfig, *config.WebServerConfig) {
	sql := "SELECT SocketServerConfig, WebServerConfig FROM serverconfig WHERE ServerGroupId = ?;"
	rows, err := dal.ChatDB().Query(sql, id)
	if err != nil {
		panic(err)
	}

	var socketServerConfig *config.SocketServerConfig
	var webServerConfig *config.WebServerConfig
	for rows.Next() {
		var socketServerConfigStr string
		var webServerConfigStr string
		err := rows.Scan(&socketServerConfigStr, &webServerConfigStr)
		if err != nil {
			panic(err)
		}

		socketServerConfig, webServerConfig = config.NewSocketServerConfig(socketServerConfigStr), config.NewWebServerConfig(webServerConfigStr)
	}

	if socketServerConfig == nil || webServerConfig == nil {
		panic(errors.New("未找到SocketServerConfig, WebServerConfig配置内容"))
	}

	return socketServerConfig, webServerConfig
}
