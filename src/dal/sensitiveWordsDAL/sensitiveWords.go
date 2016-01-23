/*
敏感词的逻辑处理包
*/
package sensitiveWordsDAL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
)

func GetList() (sensitiveWordsList []string) {
	sql := "SELECT Text FROM sensitivewords;"

	rows, err := dal.DB().Query(sql)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var text string
		err := rows.Scan(&text)
		if err != nil {
			panic(err)
		}

		sensitiveWordsList = append(sensitiveWordsList, text)
	}

	return
}
