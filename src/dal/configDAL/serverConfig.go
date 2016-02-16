package configDAL

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/ChatServer_Go/src/model/config"
)

func GetServerConfig(id int) (*config.SocketServerConfig, *config.WebServerConfig) {
	command := "SELECT SocketServerConfig, WebServerConfig FROM serverconfig WHERE ServerGroupId = ?;"

	var socketServerConfig string
	var webServerConfig string
	if err := dal.ChatDB().QueryRow(command, id).Scan(&socketServerConfig, &webServerConfig); err != nil {
		if err == sql.ErrNoRows {
			panic(errors.New("未找到SocketServerConfig, WebServerConfig配置内容"))
		} else {
			panic(errors.New(fmt.Sprintf("QueryRow失败，错误信息：%s，command:%s", err, command)))
		}
	}

	return config.NewSocketServerConfig(socketServerConfig), config.NewWebServerConfig(webServerConfig)
}
