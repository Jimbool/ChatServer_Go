package chatBLL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/clientBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/playerBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/mathUtil"
	"github.com/Jordanzuo/goutil/timeUtil"
	"time"
)

func init() {
	go displayDataSize()
}

// 显示数据大小信息(每分钟更新一次)
func displayDataSize() {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
		}
	}()

	for {
		// 刚启动时不需要显示信息，故将Sleep放在前面，而不是最后
		time.Sleep(1 * time.Minute)

		// 组装需要记录的消息
		msg := fmt.Sprintf("总共收到%s，发送%s.\t", mathUtil.GetSizeDesc(client.TotalReceiveSize()), mathUtil.GetSizeDesc(client.TotalSendSize()))
		msg += fmt.Sprintf("当前客户端数量：%d, 玩家数量：%d", clientBLL.GetClientCount(), playerBLL.GetPlayerCount())

		// 显示在控制台上，是为了便于本地使用；记录到日志文件是为了便于生产环境使用
		fmt.Println(timeUtil.Format(time.Now(), "yyyy-MM-dd HH:mm:ss:"), msg)
		logUtil.Log(msg, logUtil.Debug, true)
	}
}
