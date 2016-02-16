package configDAL

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/ChatServer_Go/src/model/config"
)

func GetAppConfig() *config.AppConfig {
	command := "SELECT AppId, AppName, AppKey FROM appconfig;"

	var appId string
	var appName string
	var appKey string
	if err := dal.ChatDB().QueryRow(command).Scan(&appId, &appName, &appKey); err != nil {
		if err == sql.ErrNoRows {
			panic(errors.New("未找到appconfig配置内容"))
		} else {
			panic(errors.New(fmt.Sprintf("QueryRow失败，错误信息：%s，command:%s", err, command)))
		}
	}

	return config.NewAppConfig(appId, appName, appKey)
}
