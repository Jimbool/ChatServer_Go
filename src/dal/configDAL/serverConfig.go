package configDAL

import (
	"errors"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/ChatServer_Go/src/model/config"
)

func GetServerConfig(id int) (*config.SocketServerConfig, *config.WebServerConfig) {
	sql := "SELECT SocketServerConfig, WebServerConfig FROM serverconfig WHERE ServerGroupId = ?;"
	rows, err := dal.ChatDB().Query(sql, id)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Query失败，错误信息：%s，sql:%s", err, sql)))
	}

	var socketServerConfig *config.SocketServerConfig
	var webServerConfig *config.WebServerConfig
	for rows.Next() {
		var socketServerConfigStr string
		var webServerConfigStr string
		err := rows.Scan(&socketServerConfigStr, &webServerConfigStr)
		if err != nil {
			panic(errors.New(fmt.Sprintf("Scan失败，错误信息：%s，sql:%s", err, sql)))
		}

		socketServerConfig = config.NewSocketServerConfig(socketServerConfigStr)
		webServerConfig = config.NewWebServerConfig(webServerConfigStr)
	}

	if socketServerConfig == nil || webServerConfig == nil {
		panic(errors.New("未找到SocketServerConfig, WebServerConfig配置内容"))
	}

	return socketServerConfig, webServerConfig
}
