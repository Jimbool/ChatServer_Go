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

	for rows.Next() {
		var socketServerConfig string
		var webServerConfig string
		err := rows.Scan(&socketServerConfig, &webServerConfig)
		if err != nil {
			panic(err)
		}

		return config.NewSocketServerConfig(socketServerConfig), config.NewWebServerConfig(webServerConfig)
	}

	panic(errors.New("未找到SocketServerConfig, WebServerConfig配置内容"))
}
