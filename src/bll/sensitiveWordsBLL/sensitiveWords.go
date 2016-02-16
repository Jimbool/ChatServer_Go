/*
屏蔽词处理包
*/
package sensitiveWordsBLL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/dal/sensitiveWordsDAL"
	"strings"
)

var (
	sensitiveWordsList = make([]string, 0, 1024)
)

func init() {
	if tmpList, err := sensitiveWordsDAL.GetList(); err != nil {
		panic(err)
	} else {
		sensitiveWordsList = tmpList
	}
}

// 重新加载
func Reload() error {
	if tmpList, err := sensitiveWordsDAL.GetList(); err != nil {
		return err
	} else {
		sensitiveWordsList = tmpList
	}

	return nil
}

// 处理屏蔽词汇
// 输入字符串
// 处理屏蔽词汇后的字符串
func HandleSensitiveWords(input string) string {
	if len(sensitiveWordsList) == 0 {
		return input
	}

	// 遍历，并将屏蔽词替换为*
	for _, item := range sensitiveWordsList {
		if strings.Contains(strings.ToUpper(input), item) {
			input = strings.Replace(input, item, "*", -1) // -1表示全部替换
		}
	}

	return input
}

// 是否包含敏感词
func IfContainsSensitiveWords(input string) bool {
	if len(sensitiveWordsList) == 0 {
		return false
	}

	// 遍历，并将屏蔽词替换为*
	for _, item := range sensitiveWordsList {
		if strings.Contains(strings.ToUpper(input), item) {
			return true
		}
	}

	return false
}
