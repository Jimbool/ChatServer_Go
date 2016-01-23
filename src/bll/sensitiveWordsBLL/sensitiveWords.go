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
	sensitiveWordsList = sensitiveWordsDAL.GetList()
}

// 重新加载
func Reload() {
	sensitiveWordsList = sensitiveWordsDAL.GetList()
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
