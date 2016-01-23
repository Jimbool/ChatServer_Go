package playerBLL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/bll/clientBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/dal/playerDAL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/ChatServer_Go/src/model/player"
	"time"
)

// 获取玩家数量
// 返回值：
// 玩家数量
func GetPlayerCount() int {
	return len(playerList)
}

// 获取所有玩家列表
// 返回值：
// 所有玩家列表
func GetPlayerList() (finalPlayerList []*player.Player) {
	for _, item := range playerList {
		finalPlayerList = append(finalPlayerList, item)
	}

	return
}

// 获取指定玩家同工会的所有玩家列表
// playerObj：指定玩家
// 返回值：
// 同工会的所有玩家列表
func GetPlayerListInSameUnion(playerObj *player.Player) (finalPlayerList []*player.Player) {
	// 筛选同一个公会的成员
	for _, item := range playerList {
		if item.UnionId == playerObj.UnionId {
			finalPlayerList = append(finalPlayerList, item)
		}
	}

	return
}

// 根据Id获取玩家对象（先从缓存中取，取不到再从数据库中去取）
// id：玩家Id
// isLoadFromDB：是否要从数据库中获取数据
// 返回值：
// 玩家对象
// 是否存在该玩家
func GetPlayer(id string, isLoadFromDB bool) (playerObj *player.Player, ok bool) {
	if id == "" {
		return nil, false
	}

	if playerObj, ok = playerList[id]; !ok {
		if isLoadFromDB {
			if playerObj, ok = playerDAL.GetPlayer(id); !ok {
				return nil, false
			}
			return playerObj, true
		}
	}

	return playerObj, true
}

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

// 根据客户端对象来断开连接
// 注销客户端连接
// 从缓存中移除玩家对象
func DisconnectByClient(clientObj *client.Client) {
	// 将玩家从缓存中移除
	if clientObj.PlayerId() != "" {
		if playerObj, ok := GetPlayer(clientObj.PlayerId(), false); ok {
			UnRegisterPlayer(playerObj)
		}
	}

	// 注销客户端连接
	clientObj.LogoutAndQuit()
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

// 更新玩家信息
// playerObj：玩家对象
// name：玩家名称
// unionId：玩家公会Id
// extraMsg：玩家透传信息
func UpdateInfo(playerObj *player.Player, name, unionId, extraMsg string) {
	playerObj.Name = name
	playerObj.UnionId = unionId
	playerObj.ExtraMsg = extraMsg

	playerDAL.UpdateInfo(playerObj)
}

// 更新登录信息
// playerObj：玩家对象
// clientObj：客户端对象
func UpdateLoginInfo(playerObj *player.Player, clientObj *client.Client) {
	playerObj.ClientId = clientObj.Id()
	playerObj.LoginTime = time.Now()

	playerDAL.UpdateLoginTime(playerObj)
}

// 更新玩家的封号状态
// playerObj：玩家对象
// isForbidden：是否封号
func UpdateForbidStatus(playerObj *player.Player, isForbidden bool) {
	playerObj.IsForbidden = isForbidden
	playerDAL.UpdateForbiddenStatus(playerObj)

	// 断开客户端连接
	if isForbidden {
		DisconnectByPlayer(playerObj)
	}
}

// 更新玩家的禁言状态
// playerObj：玩家对象
// silentEndTime：禁言结束时间
func UpdateSilentStatus(playerObj *player.Player, silentEndTime time.Time) {
	playerObj.SilentEndTime = silentEndTime
	playerDAL.UpdateSilentEndTime(playerObj)

	// 客户端退出
	if silentEndTime.Unix() > time.Now().Unix() {
		DisconnectByPlayer(playerObj)
	}
}
