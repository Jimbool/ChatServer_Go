package playerBLL

import (
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
