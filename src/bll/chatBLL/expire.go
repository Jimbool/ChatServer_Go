package chatBLL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/clientBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/playerBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/disconnectType"
	"github.com/Jordanzuo/goutil/logUtil"
	"time"
)

func init() {
	go clearExpiredClient()
}

// 清理过期的客户端
func clearExpiredClient() {
	for {
		// 休眠指定的时间（单位：秒）(放在此处是因为程序刚启动时并没有过期的客户端，所以先不用占用资源；并且此时LogPath尚未设置，如果直接执行后面的代码会出现panic异常)
		time.Sleep(configBLL.CheckExpireInterval() * time.Second)

		beforeClientCount := clientBLL.GetClientCount()
		beforePlayerCount := playerBLL.GetPlayerCount()

		// 获取过期的客户端列表
		expiredClientList := clientBLL.GetExpiredClientList()
		expiredClientCount := len(expiredClientList)
		if expiredClientCount == 0 {
			continue
		}

		for _, item := range expiredClientList {
			playerBLL.DisconnectByClient(item, disconnectType.FromExpire)
		}

		// 记录日志
		logUtil.Log(fmt.Sprintf("清理前的客户端数量为：%d，清理前的玩家数量为：%d， 本次清理不活跃的数量为：%d", beforeClientCount, beforePlayerCount, expiredClientCount), logUtil.Debug, true)
	}
}
