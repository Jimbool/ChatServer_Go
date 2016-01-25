package playerBLL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/bll/clientBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
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
// ifFromRpc：是否是来自于rpc的调用，如果是则意味着之前客户端已经关闭连接，现在需要将客户端对象从缓存中移除了；否则是客户端过期，需要关闭
func DisconnectByClient(clientObj *client.Client, isFromRpc bool) {
	// 将玩家从缓存中移除
	if clientObj.PlayerId() != "" {
		if playerObj, ok := GetPlayer(clientObj.PlayerId(), false); ok {
			UnRegisterPlayer(playerObj)
		}
	}

	// 注销客户端连接，并从缓存中移除
	if isFromRpc {
		clientBLL.UnRegisterClient(clientObj)
	} else {
		clientObj.LogoutAndQuit()
	}
}

// 根据玩家对象来断开客连接
// 注销客户端连接
// 从缓存中移除玩家对象
func DisconnectByPlayer(playerObj *player.Player) {
	// 断开客户端连接
	if playerObj.ClientId > 0 {
		if clientObj, ok := clientBLL.GetClient(playerObj.ClientId); ok {
			clientObj.LogoutAndQuit()
		}
	}

	// 将玩家从缓存中移除
	UnRegisterPlayer(playerObj)
}
