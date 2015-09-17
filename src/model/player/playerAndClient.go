package player

import (
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
)

// 玩家和客户端的映射关系类型，用于在通道中传输
type PlayerAndClient struct {
	*Player
	*client.Client
}

// 创建新的玩家和客户端映射对象
func NewPlayerAndClient(playerObj *Player, clientObj *client.Client) *PlayerAndClient {
	return &PlayerAndClient{
		playerObj,
		clientObj,
	}
}
