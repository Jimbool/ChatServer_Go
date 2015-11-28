package main

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/chatBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/rpcBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/goutil/fileUtil"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/mathUtil"
	"github.com/Jordanzuo/goutil/stringUtil"
	"github.com/Jordanzuo/goutil/timeUtil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
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

// 显示数据大小信息(每分钟更新一次)
func displayDataSize() {
	for {
		// 刚启动时不需要显示信息，故将Sleep放在前面，而不是最后
		time.Sleep(time.Minute)

		// 组装需要记录的消息
		msg := fmt.Sprintf("%s:总共收到%s，发送%s", timeUtil.Format(time.Now(), "yyyy-MM-dd HH:mm:ss"), mathUtil.GetSizeDesc(client.TotalReceiveSize()), mathUtil.GetSizeDesc(client.TotalSendSize()))
		msg += stringUtil.GetNewLineString()
		msg += fmt.Sprintf("当前客户端数量：%d, 当前玩家数量：%d", len(chatBLL.ClientList), len(chatBLL.PlayerList))

		// 显示在控制台上，是为了便于本地使用；记录到日志文件是为了便于生产环境使用
		fmt.Println(msg)
		logUtil.Log(msg, logUtil.Debug, true)
	}
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

	// 显示数据大小信息
	go displayDataSize()

	// 阻塞等待，以免main线程退出
	wg.Wait()
}
