package disconnectType

// 客户端断开连接类型
type ClientDisconnectType int

const (
	// 来自于过期检测
	FromExpire ClientDisconnectType = 1 + iota

	// 来自于Rpc检测
	FromRpc
)
