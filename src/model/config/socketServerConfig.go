package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type SocketServerConfig struct {
	ServerHost          string        // 服务器主机/IP
	ServerPort          int           // 服务器端口
	CheckExpireInterval time.Duration // 检测客户端超时的时间间隔，单位：秒
	ClientExpireSeconds time.Duration // 客户端超时的秒数
	MaxMsgLength        int           // 最大消息长度
	MaxHistoryCount     int
}

func NewSocketServerConfig(socketServerConfigStr string) *SocketServerConfig {
	var config *SocketServerConfig
	if err := json.Unmarshal([]byte(socketServerConfigStr), &config); err != nil {
		panic(errors.New(fmt.Sprintf("反序列化数据出错，原始数据位：%s，错误信息为：%s", socketServerConfigStr, err)))
	}

	return config
}

func (config *SocketServerConfig) SocketServerAddress() string {
	return fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
}
