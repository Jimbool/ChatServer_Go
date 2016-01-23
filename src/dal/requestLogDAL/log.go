/*
请求日志数据处理包
*/
package requestLogDAL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"time"
)

func Insert(apiName, content string) {
	sql := "INSERT INTO request_log(APIName, Content, Crdate) VALUES(?, ?, ?);"
	stmt, err := dal.DB().Prepare(sql)
	if err != nil {
		panic(err)
	}

	// 最后关闭
	defer stmt.Close()

	_, err = stmt.Exec(apiName, content, time.Now())
	if err != nil {
		panic(err)
	}
}
