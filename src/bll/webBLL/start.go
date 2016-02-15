package webBLL

import (
	"errors"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"net/http"
)

func init() {
	go start()
}

func start() {
	// 设置访问的路由
	mux := new(SelfDefineMux)

	// 启动Web服务器监听
	err := http.ListenAndServe(configBLL.WebServerAddress(), mux)
	if err != nil {
		panic(errors.New(fmt.Sprintf("ListenAndServe失败，错误信息为：%s", err)))
	}
}
