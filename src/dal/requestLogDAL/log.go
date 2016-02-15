/*
请求日志数据处理包
*/
package requestLogDAL

import (
	"errors"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"time"
)

func Insert(apiName, content string) error {
	sql := "INSERT INTO request_log(APIName, ServerGroupId, Content, Crdate) VALUES(?, ?, ?, ?);"
	stmt, err := dal.ChatDB().Prepare(sql)
	if err != nil {
		return errors.New(fmt.Sprintf("Prepare失败，错误信息：%s，sql:%s", err, sql))
	}

	// 最后关闭
	defer stmt.Close()

	_, err = stmt.Exec(apiName, dal.ServerGroupId, content, time.Now())
	if err != nil {
		return errors.New(fmt.Sprintf("Exec失败，错误信息：%s，sql:%s", err, sql))
	}

	return nil
}
