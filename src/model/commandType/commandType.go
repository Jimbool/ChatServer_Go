/*
聊天命令类型包，定义了聊天频道发送消息的命令类型
*/
package commandType

// 定义聊天命令类型对象
type CommandType int

const (
	// 登陆
	Login CommandType = 1 + iota

	// 登出
	Logout

	// 发送消息
	SendMessage

	// 更新玩家信息
	UpdatePlayerInfo
)
