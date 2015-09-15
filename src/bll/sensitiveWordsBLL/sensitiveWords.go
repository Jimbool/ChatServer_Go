package sensitiveWordsBLL

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const (
	// 存放屏蔽词库的文件名称
	SENSITIVEWORDS_FILENAME = "SensitiveWords.txt"
)

var (
	// 屏蔽词列表
	SensitiveWordsList = make([]string, 0, 1024)
)

func init() {
	//打开文件
	file, err := os.Open(SENSITIVEWORDS_FILENAME)
	if err != nil {
		// 由于屏蔽词库不是在程序启动时加载的，所以即便是失败也要不影响程序的进行，所以此处不用panic
		return
	}
	defer file.Close()

	//读取文件
	buf := bufio.NewReader(file)
	for {
		//按行读取
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}

		//将byte[]转换为string
		lineStr := string(line)

		// 添加到列表中
		SensitiveWordsList = append(SensitiveWordsList, lineStr)
	}
}

// 处理屏蔽词汇
// 输入字符串
// 处理屏蔽词汇后的字符串
func HandleSensitiveWords(input string) string {
	if len(SensitiveWordsList) == 0 {
		return input
	}

	// 遍历，并将屏蔽词替换为*
	for _, item := range SensitiveWordsList {
		if strings.Contains(strings.ToUpper(input), item) {
			input = strings.Replace(input, item, "*", -1) // -1表示全部替换
		}
	}

	return input
}
