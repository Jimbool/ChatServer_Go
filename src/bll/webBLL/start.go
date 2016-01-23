package webBLL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"github.com/Jordanzuo/goutil/logUtil"
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
		msg := fmt.Sprintf("ListenAndServe失败，错误信息为：%s", err)
		fmt.Println(msg)
		logUtil.Log(msg, logUtil.Error, true)
	}
}
