/*
敏感词的逻辑处理包
*/
package sensitiveWordsDAL

import (
	"errors"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
)

func GetList() (sensitiveWordsList []string) {
	sql := "SELECT Words FROM b_sensitive_words_c;"
	rows, err := dal.ModelDB().Query(sql)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Query失败，错误信息：%s，sql:%s", err, sql)))
	}

	for rows.Next() {
		var words string
		err := rows.Scan(&words)
		if err != nil {
			panic(errors.New(fmt.Sprintf("Scan失败，错误信息：%s，sql:%s", err, sql)))
		}

		sensitiveWordsList = append(sensitiveWordsList, words)
	}

	return
}
