/*
敏感词的逻辑处理包
*/
package sensitiveWordsDAL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/goutil/logUtil"
	"time"
)

func GetList() (sensitiveWordsList []string) {
	sql := "SELECT Words FROM b_sensitive_words_c;"

	// log begin
	start := time.Now().Unix()
	defer func() {
		end := time.Now().Unix()
		logUtil.Log(fmt.Sprintf("%s执行耗时%d", sql, (end-start)), logUtil.Fatal, true)
	}()
	// log end

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
