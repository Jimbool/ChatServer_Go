package playerBLL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/dal/playerDAL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/player"
	"sync"
)

var (
	// 玩家列表
	playerList = make(map[string]*player.Player, 1024)

	// 读写锁
	mutex sync.RWMutex
)

// 注册玩家对象到缓存中
// playerObj：玩家对象
func RegisterPlayer(playerObj *player.Player) {
	mutex.Lock()
	defer mutex.Unlock()

	playerList[playerObj.Id] = playerObj
}

// 从缓存中取消玩家注册
// playerObj：玩家对象
func UnRegisterPlayer(playerObj *player.Player) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(playerList, playerObj.Id)
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

	mutex.RLock()
	if playerObj, ok = playerList[id]; !ok {
		mutex.RUnlock()
		if isLoadFromDB {
			if playerObj, ok = playerDAL.GetPlayer(id); !ok {
				return nil, false
			} else {
				return playerObj, true
			}
		} else {
			return nil, false
		}
	} else {
		mutex.RUnlock()
		return playerObj, true
	}
}

// 获取玩家数量
// 返回值：
// 玩家数量
func GetPlayerCount() int {
	mutex.RLock()
	defer mutex.RUnlock()

	return len(playerList)
}

// 获取所有玩家列表
// 返回值：
// 所有玩家列表
func GetPlayerList() (finalPlayerList []*player.Player) {
	mutex.RLock()
	defer mutex.RUnlock()

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
	mutex.RLock()
	defer mutex.RUnlock()

	for _, item := range playerList {
		if item.UnionId == playerObj.UnionId {
			finalPlayerList = append(finalPlayerList, item)
		}
	}

	return
}
