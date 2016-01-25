package playerBLL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/bll/clientBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/ChatServer_Go/src/model/disconnectType"
	"github.com/Jordanzuo/ChatServer_Go/src/model/player"
	"github.com/Jordanzuo/goutil/logUtil"
	"time"
)

var (
	// 玩家列表
	playerList = make(map[string]*player.Player)

	playerAddChan    = make(chan *player.Player, 50)
	playerRemoveChan = make(chan *player.Player, 50)
)

func init() {
	go handlePlayerChannel()
}

// 处理增加、删除客户端channel；增加、删除玩家的channel的逻辑
func handlePlayerChannel() {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
		}
	}()

	for {
		select {
		case playerObj := <-playerAddChan:
			addPlayer(playerObj)
		case playerObj := <-playerRemoveChan:
			removePlayer(playerObj)
		default:
			// 休眠一下，防止CPU过高
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func addPlayer(playerObj *player.Player) {
	playerList[playerObj.Id] = playerObj
}

func removePlayer(playerObj *player.Player) {
	delete(playerList, playerObj.Id)
}

// 注册玩家对象到缓存中
// playerObj：玩家对象
func RegisterPlayer(playerObj *player.Player) {
	playerAddChan <- playerObj
}

// 从缓存中取消玩家注册
// playerObj：玩家对象
func UnRegisterPlayer(playerObj *player.Player) {
	playerRemoveChan <- playerObj
}

// 从缓存中取消玩家Id对应的玩家
// id：玩家Id
func UnRegisterPlayerId(id string) {
	if playerObj, ok := GetPlayer(id, false); ok {
		UnRegisterPlayer(playerObj)
	}
}

// 根据客户端对象来断开连接
// 注销客户端连接
// 从缓存中移除玩家对象
// clientObj：客户端对象
// clientDisconnectType：客户端断开连接的类型
func DisconnectByClient(clientObj *client.Client, clientDisconnectType disconnectType.ClientDisconnectType) {
	// 将玩家从缓存中移除
	if clientObj.PlayerId() != "" {
		if playerObj, ok := GetPlayer(clientObj.PlayerId(), false); ok {
			UnRegisterPlayer(playerObj)
		}
	}

	// 如果是来自于客户端过期，则将客户端登出并断开连接；
	// 如果是来自于RPC，则将客户端从缓存中移除
	switch clientDisconnectType {
	case disconnectType.FromExpire:
		clientObj.LogoutAndQuit()
	case disconnectType.FromRpc:
		clientBLL.UnRegisterClient(clientObj)
	}
}

// 根据玩家对象来断开客连接
// 注销客户端连接
// 从缓存中移除玩家对象
// playerObj：玩家对象
// playerDisconnectType：玩家断开连接的类型
func DisconnectByPlayer(playerObj *player.Player, playerDisconnectType disconnectType.PlayerDisconnectType) {
	if playerObj.ClientId > 0 {
		if clientObj, ok := clientBLL.GetClient(playerObj.ClientId); ok {
			// 先发送指定类型的消息
			switch playerDisconnectType {
			case disconnectType.FromForbid:
				SendForbidMsg(clientObj)
			case disconnectType.FromSilent:
				SendSilentMsg(clientObj)
			}

			// 最后取消与玩家的关联，并断开客户端连接
			clientObj.LogoutAndQuit()
		}
	}

	// 将玩家从缓存中移除
	UnRegisterPlayer(playerObj)
}
