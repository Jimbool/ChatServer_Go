/*
敏感词的逻辑处理包
*/
package sensitiveWordsDAL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/goutil/logUtil"
)

func GetList() (sensitiveWordsList []string, err error) {
	command := "SELECT Words FROM b_sensitive_words_c;"
	rows, err := dal.ModelDB().Query(command)
	if err != nil {
		logUtil.Log(fmt.Sprintf("Query失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return
	}

	// 最后关闭
	defer rows.Close()

	for rows.Next() {
		var words string
		if err = rows.Scan(&words); err != nil {
			logUtil.Log(fmt.Sprintf("Scan失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
			return
		}

		sensitiveWordsList = append(sensitiveWordsList, words)
	}

	return
}
