/*
封装数据库对象的包
*/
package dal

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	chatDB  *sql.DB
	modelDB *sql.DB
	gameDB  *sql.DB
)

// 初始化数据库连接相关的配置
func init() {
	chatDB = openMysqlConnection(ChatDBConnection, ChatDBMaxOpenConns, ChatDBMaxIdleConns)
	modelDB = openMysqlConnection(ModelDBConnection, ModelDBMaxOpenConns, ModelDBMaxIdleConns)
	gameDB = openMysqlConnection(GameDBConnection, GameDBMaxOpenConns, GameDBMaxIdleConns)

	// 启动一个Goroutine一直ping数据库，以免被数据库认为过期而关掉
	go ping()
}

// 获取聊天数据库对象
// 返回值：
// 聊天数据库对象
func ChatDB() *sql.DB {
	return chatDB
}

// 获取模型数据库对象
// 返回值：
// 模型数据库对象
func ModelDB() *sql.DB {
	return modelDB
}

// 获取游戏数据库对象
// 返回值：
// 游戏数据库对象
func GameDB() *sql.DB {
	return gameDB
}

func openMysqlConnection(connectionString string, maxOpenConns, maxIdleConns int) *sql.DB {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(errors.New(fmt.Sprintf("打开数据库失败,连接字符串为：%s", connectionString)))
	}

	if maxOpenConns > 0 && maxIdleConns > 0 {
		db.SetMaxOpenConns(maxOpenConns)
		db.SetMaxIdleConns(maxIdleConns)
	}

	if err := db.Ping(); err != nil {
		panic(errors.New(fmt.Sprintf("Ping数据库失败,连接字符串为：%s", connectionString)))
	}

	return db
}

// 每分钟ping一次数据库
func ping() {
	for {
		time.Sleep(time.Minute)

		chatDB.Ping()
		gameDB.Ping()
	}
}
