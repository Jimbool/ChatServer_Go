/*
聊天频道类型包，定义了聊天频道的类型，包括世界、公会、私聊等3个频道
*/
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
