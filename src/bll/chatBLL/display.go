package chatBLL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/clientBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/playerBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/mathUtil"
	"time"
)

func init() {
	go displayDataSize()
}

// 显示数据大小信息(每5分钟更新一次)
func displayDataSize() {
	for {
		// 刚启动时不需要显示信息，故将Sleep放在前面，而不是最后
		time.Sleep(5 * time.Minute)

		// 组装需要记录的信息
		msg := fmt.Sprintf("总共收到%s，发送%s.\t", mathUtil.GetSizeDesc(client.TotalReceiveSize()), mathUtil.GetSizeDesc(client.TotalSendSize()))
		msg += fmt.Sprintf("当前客户端数量：%d, 玩家数量：%d", clientBLL.GetClientCount(), playerBLL.GetPlayerCount())
		logUtil.Log(msg, logUtil.Debug, true)
	}
}
