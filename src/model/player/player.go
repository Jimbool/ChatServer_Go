package player

// 定义玩家对象
type Player struct {
	// 玩家Id
	Id string

	// 玩家名称
	Name string

	// 玩家公会Id
	UnionId string

	// 玩家额外信息
	ExtraMsg interface{}
}

// 构造新的Plyaer对象
// id：玩家Id
// name：玩家名称
// unionId：玩家所在公会的Id
// extraData；额外信息
func NewPlayer(id, name, unionId string, extraMsg interface{}) *Player {
	return &Player{
		Id:       id,
		Name:     name,
		UnionId:  unionId,
		ExtraMsg: extraMsg,
	}
}

// 更新玩家信息
// name：玩家名称
// unionId：玩家所在公会的Id
// extraData；额外信息
func (p *Player) Update(name, unionId string, extraMsg interface{}) {
	p.Name = name
	p.UnionId = unionId
	p.ExtraMsg = extraMsg
}
