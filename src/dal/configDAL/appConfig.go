package configDAL

import (
	"errors"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/ChatServer_Go/src/model/config"
)

func GetAppConfig() *config.AppConfig {
	sql := "SELECT AppId, AppName, AppKey FROM appconfig;"
	rows, err := dal.ChatDB().Query(sql)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var appId string
		var appName string
		var appKey string
		err := rows.Scan(&appId, &appName, &appKey)
		if err != nil {
			panic(err)
		}

		return config.NewAppConfig(appId, appName, appKey)
	}

	panic(errors.New("未找到config配置内容"))
}
