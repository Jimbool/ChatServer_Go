package playerBLL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/dal/playerDAL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/ChatServer_Go/src/model/disconnectType"
	"github.com/Jordanzuo/ChatServer_Go/src/model/player"
	"time"
)

// 注册新玩家
// id：玩家Id
// name：玩家名称
// unionId：玩家公会Id
// extraMsg：玩家透传信息
// 返回值：
// 玩家对象
func RegisterNewPlayer(id, name, unionId, extraMsg string) *player.Player {
	playerObj := player.InitPlayer(id, name, unionId, extraMsg)
	playerDAL.Insert(playerObj)

	return playerObj
}

// 更新玩家信息
// playerObj：玩家对象
// name：玩家名称
// unionId：玩家公会Id
// extraMsg：玩家透传信息
func UpdateInfo(playerObj *player.Player, name, unionId, extraMsg string) {
	playerObj.Name = name
	playerObj.UnionId = unionId
	playerObj.ExtraMsg = extraMsg

	go playerDAL.UpdateInfo(playerObj)
}

// 更新登录信息
// playerObj：玩家对象
// clientObj：客户端对象
// isNewPlayer：是否是新玩家
func UpdateLoginInfo(playerObj *player.Player, clientObj *client.Client, isNewPlayer bool) {
	playerObj.ClientId = clientObj.Id()
	playerObj.LoginTime = time.Now()

	// 如果不是新玩家则更新登录时间，否则使用创建时指定的登录时间
	if !isNewPlayer {
		go playerDAL.UpdateLoginTime(playerObj)
	}
}

// 更新玩家的封号状态
// playerObj：玩家对象
// isForbidden：是否封号
func UpdateForbidStatus(playerObj *player.Player, isForbidden bool) {
	playerObj.IsForbidden = isForbidden
	playerDAL.UpdateForbiddenStatus(playerObj)

	// 断开客户端连接
	if isForbidden {
		disconnectByPlayer(playerObj, disconnectType.FromForbid)
	}
}

// 更新玩家的禁言状态
// playerObj：玩家对象
// silentEndTime：禁言结束时间
func UpdateSilentStatus(playerObj *player.Player, silentEndTime time.Time) {
	playerObj.SilentEndTime = silentEndTime
	playerDAL.UpdateSilentEndTime(playerObj)
}

// 获取游戏玩家名称
// playerObj：玩家对象
// 返回值：
// 玩家名称
// 玩家公会Id
// 是否存在玩家
func GetGamePlayer(id string) (string, string, bool) {
	return playerDAL.GetGamePlayer(id)
}
