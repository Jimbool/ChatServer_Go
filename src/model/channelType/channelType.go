package channelType

// 聊天频道类型定义
type ChannelType int

const (
	// 世界频道
	World ChannelType = 1 + iota

	// 公会频道
	Union

	// 私聊频道
	Private
)
