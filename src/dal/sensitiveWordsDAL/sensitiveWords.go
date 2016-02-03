/*
敏感词的逻辑处理包
*/
package sensitiveWordsDAL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
)

func GetList() (sensitiveWordsList []string) {
	sql := "SELECT Words FROM b_sensitive_words_c;"
	rows, err := dal.ModelDB().Query(sql)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var words string
		err := rows.Scan(&words)
		if err != nil {
			panic(err)
		}

		sensitiveWordsList = append(sensitiveWordsList, words)
	}

	return
}
