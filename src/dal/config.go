/*
数据库配置逻辑处理包
*/
package dal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

const (
	// 配置文件名称
	CONFIG_FILE_NAME = "config.ini"
)

var (
	// 服务器组Id
	ServerGroupId int

	// 聊天数据库连接字符串
	ChatDBConnection string

	// 模型数据库连接字符串
	ModelDBConnection string

	// 游戏数据库连接字符串
	GameDBConnection string

	// 数据库的最大连接数
	MaxOpenConns int

	// 数据库的最大空闲数
	MaxIdleConns int
)

func init() {
	// 由于服务器的运行依赖于init中执行的逻辑，所以如果出现任何的错误都直接panic，让程序启动失败；而不是让它启动成功，但是在运行时出现错误

	// 读取配置文件（一次性读取整个文件，则使用ioutil）
	bytes, err := ioutil.ReadFile(CONFIG_FILE_NAME)
	if err != nil {
		panic(err)
	}

	// 使用json反序列化
	config := make(map[string]interface{})
	if err = json.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}

	// 解析ServerGroupId
	serverGroupId, ok := config["ServerGroupId"]
	if !ok {
		panic(errors.New("不存在名为ServerGroupId的配置或配置为空"))
	}
	serverGroupId_float, ok := serverGroupId.(float64)
	if !ok {
		panic(errors.New("ServerGroupId必须为int型"))
	}

	// 设置ServerGroupId
	ServerGroupId = int(serverGroupId_float)

	// 解析ChatDBConnection
	chatDBConnection, ok := config["ChatDBConnection"]
	if !ok {
		panic(errors.New("不存在名为ChatDBConnection的配置或配置为空"))
	}
	chatDBConnection_string, ok := chatDBConnection.(string)
	if !ok {
		panic(errors.New("ChatDBConnection必须为字符串类型"))
	}

	// 设置ChatDBConnection
	ChatDBConnection = chatDBConnection_string

	// 解析ModelDBConnection
	modelDBConnection, ok := config["ModelDBConnection"]
	if !ok {
		panic(errors.New("不存在名为ModelDBConnection的配置或配置为空"))
	}
	modelDBConnection_string, ok := modelDBConnection.(string)
	if !ok {
		panic(errors.New("ModelDBConnection必须为字符串类型"))
	}

	// 设置ModelDBConnection
	ModelDBConnection = modelDBConnection_string

	// 解析GameDBConnection
	gameDBConnection, ok := config["GameDBConnection"]
	if !ok {
		panic(errors.New("不存在名为GameDBConnection的配置或配置为空"))
	}
	gameDBConnection_string, ok := gameDBConnection.(string)
	if !ok {
		panic(errors.New("GameDBConnection必须为字符串类型"))
	}

	// 设置GameDBConnection
	GameDBConnection = gameDBConnection_string

	// 解析SERVER_PORT
	maxOpenConns, ok := config["MaxOpenConns"]
	if !ok {
		panic(errors.New("不存在名为MaxOpenConns的配置或配置为空"))
	}
	maxOpenConns_float, ok := maxOpenConns.(float64)
	if !ok {
		panic(errors.New("MaxOpenConns必须是int型"))
	}

	// 设置MaxOpenConns
	MaxOpenConns = int(maxOpenConns_float)

	// 解析MaxIdleConns
	maxIdleConns, ok := config["MaxIdleConns"]
	if !ok {
		panic(errors.New("不存在名为MaxIdleConns的配置或配置为空"))
	}
	maxIdleConns_float, ok := maxIdleConns.(float64)
	if !ok {
		panic(errors.New("MaxIdleConns必须是int型"))
	}

	// 设置MaxIdleConns
	MaxIdleConns = int(maxIdleConns_float)
}
