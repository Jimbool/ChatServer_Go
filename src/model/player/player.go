/*
玩家对象包，定义了玩家对象
*/
package player

import (
	"time"
)

// 定义玩家对象
type Player struct {
	// 玩家Id
	Id string

	// 玩家名称
	Name string

	// 玩家公会Id
	UnionId string

	// 额外透传信息
	ExtraMsg string

	//注册时间
	RegisterTime time.Time `json:-`

	//登录时间
	LoginTime time.Time `json:-`

	//是否封号
	IsForbidden bool `json:-`

	//禁言结束时间
	SilentEndTime time.Time `json:-`

	// 客户端Id
	ClientId int32 `json:-`
}

// 初始化一个玩家对象
func InitPlayer(Id, Name, UnionId string, ExtraMsg string) *Player {
	return &Player{
		Id:            Id,
		Name:          Name,
		UnionId:       UnionId,
		ExtraMsg:      ExtraMsg,
		RegisterTime:  time.Now(),
		LoginTime:     time.Now(),
		IsForbidden:   false,
		SilentEndTime: time.Now(),
		ClientId:      0,
	}
}

// 使用现有数据构造一个新的玩家对象
func NewPlayer(Id, Name, UnionId string, ExtraMsg string, registerTime, loginTime time.Time, isForbidden bool, silentEndTime time.Time) *Player {
	return &Player{
		Id:            Id,
		Name:          Name,
		UnionId:       UnionId,
		ExtraMsg:      ExtraMsg,
		RegisterTime:  registerTime,
		LoginTime:     loginTime,
		IsForbidden:   isForbidden,
		SilentEndTime: silentEndTime,
		ClientId:      0,
	}
}

// 判断玩家是否处于禁言状态
// 返回值：
// 是否处于禁言状态
// 禁言剩余分钟数
func (playerObj *Player) IsInSilent() (bool, int) {
	leftSeconds := playerObj.SilentEndTime.Unix() - time.Now().Unix()
	if leftSeconds <= 0 {
		return false, 0
	} else {
		if leftSeconds%60 == 0 {
			return true, int(leftSeconds / 60)
		} else {
			return true, int(leftSeconds/60) + 1
		}
	}
}
