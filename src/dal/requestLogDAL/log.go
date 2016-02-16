/*
请求日志数据处理包
*/
package requestLogDAL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/goutil/logUtil"
	"time"
)

func Insert(apiName, content string) error {
	command := "INSERT INTO request_log(APIName, ServerGroupId, Content, Crdate) VALUES(?, ?, ?, ?);"
	stmt, err := dal.ChatDB().Prepare(command)
	if err != nil {
		logUtil.Log(fmt.Sprintf("Prepare失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(apiName, dal.ServerGroupId, content, time.Now()); err != nil {
		logUtil.Log(fmt.Sprintf("Exec失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	return nil
}
