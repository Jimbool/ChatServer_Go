package configBLL

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

const (
	// 配置文件名称
	CONFIG_FILE_NAME = "config.ini"
)

var (
	// 服务器监听地址
	ServerAddress string

	// 服务端检测过期客户端的时间间隔，单位：秒
	CheckExpiredInterval time.Duration

	// 客户端过期的秒数
	ClientExpiredSeconds time.Duration

	// 登陆key
	LoginKey string

	// 最大消息的长度
	MaxMsgLength int
)

func init() {
	// 由于服务器的运行依赖于init中执行的逻辑，所以如果出现任何的错误都直接panic，让程序启动失败；而不是让它启动成功，但是在运行时出现错误

	// 读取配置文件（一次性读取整个文件，则使用ioutil）
	bytes, err := ioutil.ReadFile(CONFIG_FILE_NAME)
	if err != nil {
		panic(err)
	}

	// 使用json反序列化
	config := make(map[string]interface{})
	if err = json.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}

	// 解析SERVER_HOST
	serverHost, ok := config["SERVER_HOST"]
	if !ok {
		panic(errors.New("不存在名为SERVER_HOST的配置或配置为空"))
	}
	serverHost_string, ok := serverHost.(string)
	if !ok {
		panic(errors.New("SERVER_HOST必须是字符串类型"))
	}

	// SERVER_PORT
	serverPort, ok := config["SERVER_PORT"]
	if !ok {
		panic(errors.New("不存在名为SERVER_PORT的配置或配置为空"))
	}
	serverPort_int, ok := serverPort.(float64)
	if !ok {
		panic(errors.New("SERVER_PORT必须是int型"))
	}

	// 设置ServerAddress
	ServerAddress = fmt.Sprintf("%s:%d", serverHost_string, int(serverPort_int))

	// 解析CHECK_EXPIRED_INTERVAL
	checkExpiredInterval, ok := config["CHECK_EXPIRED_INTERVAL"]
	if !ok {
		panic(errors.New("不存在名为CHECK_EXPIRED_INTERVAL的配置或配置为空"))
	}
	checkExpiredInterval_int, ok := checkExpiredInterval.(float64)
	if !ok {
		panic(errors.New("CHECK_EXPIRED_INTERVAL必须是int型"))
	}

	// 设置CheckExpiredInterval参数
	CheckExpiredInterval = time.Duration(int(checkExpiredInterval_int))

	// 解析ClientExpiredSeconds
	clientExpiredSeconds, ok := config["CLIENT_EXPIRED_SECONDS"]
	if !ok {
		panic(errors.New("不存在名为CLIENT_EXPIRED_SECONDS的配置或配置为空"))
	}
	clientExpiredSeconds_int, ok := clientExpiredSeconds.(float64)
	if !ok {
		panic(errors.New("CLIENT_EXPIRED_SECONDS必须是int型"))
	}

	// 设置client的参数：ClientExpiredSeconds
	ClientExpiredSeconds = time.Duration(int(clientExpiredSeconds_int))

	// 解析LOGIN_KEY
	loginKey, ok := config["LOGIN_KEY"]
	if !ok {
		panic(errors.New("不存在名为LOGIN_KEY的配置或配置为空"))
	}
	loginKey_string, ok := loginKey.(string)
	if !ok {
		panic(errors.New("LOGIN_KEY必须是string型"))
	}

	// 设置参数LoginKey
	LoginKey = loginKey_string

	// 解析MAX_MSG_LENGTH
	maxMsgLength, ok := config["MAX_MSG_LENGTH"]
	if !ok {
		panic(errors.New("不存在名为MAX_MSG_LENGTH的配置或配置为空"))
	}
	maxMsgLength_int, ok := maxMsgLength.(float64)
	if !ok {
		panic(errors.New("MAX_MSG_LENGTH必须是int型"))
	}

	// 设置参数MaxMsgLength
	MaxMsgLength = int(maxMsgLength_int)
}
