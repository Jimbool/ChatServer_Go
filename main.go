package main

import (
	"github.com/Jordanzuo/ChatServer_Go/src/bll/rpcBLL"
	"github.com/Jordanzuo/goutil/fileUtil"
	"github.com/Jordanzuo/goutil/logUtil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

import (
	_ "github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	_ "github.com/Jordanzuo/ChatServer_Go/src/bll/sensitiveWordsBLL"
	_ "github.com/Jordanzuo/ChatServer_Go/src/bll/webBLL"
)

const (
	// 日志文件路径后缀
	LOG_PATH_SUFFIX = "LOG"
)

var (
	wg sync.WaitGroup
)

func init() {
	// 设置日志文件的存储目录
	logUtil.SetLogPath(filepath.Join(fileUtil.GetCurrentPath(), LOG_PATH_SUFFIX))

	// 设置WaitGroup需要等待的数量
	wg.Add(1)
}

// 处理系统信号
func signalProc() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	for {
		// 准备接收信息
		<-sigs

		// 一旦收到信号，则表明管理员希望退出程序，则先保存信息，然后退出
		os.Exit(0)
	}
}

func main() {
	// 处理系统信号
	go signalProc()

	// 启动服务器
	go rpcBLL.StartServer(&wg)

	// 阻塞等待，以免main线程退出
	wg.Wait()
}
