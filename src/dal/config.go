/*
数据库配置逻辑处理包
*/
package dal

import (
	"encoding/json"
	"errors"
	"fmt"
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

	// 聊天数据库的最大连接数
	ChatDBMaxOpenConns int

	// 聊天数据库的最大空闲数
	ChatDBMaxIdleConns int

	// 模型数据库连接字符串
	ModelDBConnection string

	// 模型数据库的最大连接数
	ModelDBMaxOpenConns int

	// 模型数据库的最大空闲数
	ModelDBMaxIdleConns int

	// 游戏数据库连接字符串
	GameDBConnection string

	// 游戏数据库的最大连接数
	GameDBMaxOpenConns int

	// 游戏数据库的最大空闲数
	GameDBMaxIdleConns int
)

func init() {
	// 由于服务器的运行依赖于init中执行的逻辑，所以如果出现任何的错误都直接panic，让程序启动失败；而不是让它启动成功，但是在运行时出现错误

	// 读取配置文件（一次性读取整个文件，则使用ioutil）
	bytes, err := ioutil.ReadFile(CONFIG_FILE_NAME)
	if err != nil {
		panic(errors.New("读取配置文件的内容出错"))
	}

	// 使用json反序列化
	config := make(map[string]interface{})
	if err = json.Unmarshal(bytes, &config); err != nil {
		panic(errors.New("反序列化配置文件的内容出错"))
	}

	// 解析ServerGroupId
	ServerGroupId = initServerGroup(config)

	// 解析ChatDBConnection配置
	ChatDBConnection, ChatDBMaxOpenConns, ChatDBMaxIdleConns = initDBConnection(config, "ChatDBConnection", "ChatDBMaxOpenConns", "ChatDBMaxIdleConns")

	// 解析ModelDBConnection配置
	ModelDBConnection, ModelDBMaxOpenConns, ModelDBMaxIdleConns = initDBConnection(config, "ModelDBConnection", "ModelDBMaxOpenConns", "ModelDBMaxIdleConns")

	// 解析GameDBConnection配置
	GameDBConnection, GameDBMaxOpenConns, GameDBMaxIdleConns = initDBConnection(config, "GameDBConnection", "GameDBMaxOpenConns", "GameDBMaxIdleConns")
}

func initServerGroup(config map[string]interface{}) int {
	// 解析ServerGroupId
	serverGroupId, ok := config["ServerGroupId"]
	if !ok {
		panic(errors.New("不存在名为ServerGroupId的配置或配置为空"))
	}
	serverGroupId_float, ok := serverGroupId.(float64)
	if !ok {
		panic(errors.New("ServerGroupId必须为int型"))
	}

	return int(serverGroupId_float)
}

func initDBConnection(config map[string]interface{}, dbConnectionName, maxOpenConnsName, maxIdleConnsName string) (string, int, int) {
	// 解析DBConnection
	dbConnection, ok := config[dbConnectionName]
	if !ok {
		panic(errors.New(fmt.Sprintf("不存在名为%s的配置或配置为空", dbConnectionName)))
	}
	dbConnection_string, ok := dbConnection.(string)
	if !ok {
		panic(errors.New(fmt.Sprintf("%s必须为字符串类型", dbConnectionName)))
	}

	// 解析MaxOpenConns
	maxOpenConns, ok := config[maxOpenConnsName]
	if !ok {
		panic(errors.New(fmt.Sprintf("不存在名为%s的配置或配置为空", maxOpenConnsName)))
	}
	maxOpenConns_float, ok := maxOpenConns.(float64)
	if !ok {
		panic(errors.New(fmt.Sprintf("%s必须是int型", maxOpenConnsName)))
	}

	// 解析MaxIdleConns
	maxIdleConns, ok := config[maxIdleConnsName]
	if !ok {
		panic(errors.New(fmt.Sprintf("不存在名为%s的配置或配置为空", maxIdleConnsName)))
	}
	maxIdleConns_float, ok := maxIdleConns.(float64)
	if !ok {
		panic(errors.New(fmt.Sprintf("%s必须是int型", maxIdleConnsName)))
	}

	return dbConnection_string, int(maxOpenConns_float), int(maxIdleConns_float)
}
