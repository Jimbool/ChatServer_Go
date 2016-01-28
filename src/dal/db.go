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
	db *sql.DB
)

// 初始化数据库连接相关的配置
func init() {
	db = openMysqlConnection(DBConnection, MaxOpenConns, MaxIdleConns)

	// 启动一个Goroutine一直ping数据库，以免被数据库认为过期而关掉
	go ping()
}

// 获取数据库对象
// 返回值：
// 数据库对象
func DB() *sql.DB {
	return db
}

func openMysqlConnection(connectionString string, maxOpenConns, maxIdleConns int) *sql.DB {
	// 建立数据库连接
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(errors.New(fmt.Sprintf("打开游戏数据库失败,连接字符串为：%s", connectionString)))
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.Ping()

	return db
}

func ping() {
	// 每5秒ping一次数据库
	for {
		time.Sleep(5 * time.Second)
		db.Ping()
	}
}
