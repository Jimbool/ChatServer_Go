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
// 错误对象
func RegisterNewPlayer(id, name, unionId, extraMsg string) (*player.Player, error) {
	playerObj := player.InitPlayer(id, name, unionId, extraMsg)
	if err := playerDAL.Insert(playerObj); err != nil {
		return nil, err
	}

	return playerObj, nil
}

// 更新玩家信息
// playerObj：玩家对象
// name：玩家名称
// unionId：玩家公会Id
// extraMsg：玩家透传信息
func UpdateInfo(playerObj *player.Player, name, unionId, extraMsg string) error {
	playerObj.Name = name
	playerObj.UnionId = unionId
	playerObj.ExtraMsg = extraMsg

	return playerDAL.UpdateInfo(playerObj)
}

// 更新登录信息
// playerObj：玩家对象
// clientObj：客户端对象
// isNewPlayer：是否是新玩家
func UpdateLoginInfo(playerObj *player.Player, clientObj *client.Client, isNewPlayer bool) error {
	playerObj.ClientId = clientObj.Id()
	playerObj.LoginTime = time.Now()

	// 如果不是新玩家则更新登录时间，否则使用创建时指定的登录时间
	if !isNewPlayer {
		if err := playerDAL.UpdateLoginTime(playerObj); err != nil {
			return err
		}
	}

	return nil
}

// 更新玩家的封号状态
// playerObj：玩家对象
// isForbidden：是否封号
func UpdateForbidStatus(playerObj *player.Player, isForbidden bool) error {
	playerObj.IsForbidden = isForbidden
	if err := playerDAL.UpdateForbiddenStatus(playerObj); err != nil {
		return err
	}

	// 断开客户端连接
	if isForbidden {
		disconnectByPlayer(playerObj, disconnectType.FromForbid)
	}

	return nil
}

// 更新玩家的禁言状态
// playerObj：玩家对象
// silentEndTime：禁言结束时间
func UpdateSilentStatus(playerObj *player.Player, silentEndTime time.Time) error {
	playerObj.SilentEndTime = silentEndTime
	return playerDAL.UpdateSilentEndTime(playerObj)
}

// 获取游戏玩家名称
// playerObj：玩家对象
// 返回值：
// 玩家名称
// 玩家公会Id
// 是否存在玩家
// 错误对象
func GetGamePlayer(id string) (string, string, bool, error) {
	return playerDAL.GetGamePlayer(id)
}
