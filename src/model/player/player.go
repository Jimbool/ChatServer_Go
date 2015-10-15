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

	// 客户端Id
	clientId int32
}

// 构造新的Plyaer对象
// Id：玩家Id
// Name：玩家名称
// UnionId：玩家所在公会的Id
// extraData；额外信息
func NewPlayer(Id, Name, UnionId string, ExtraMsg interface{}, clientId int32) *Player {
	return &Player{
		Id:       Id,
		Name:     Name,
		UnionId:  UnionId,
		ExtraMsg: ExtraMsg,
		clientId: clientId,
	}
}

// 更新玩家信息
// Name：玩家名称
// UnionId：玩家所在公会的Id
// extraData；额外信息
func (p *Player) Update(Name, UnionId string, ExtraMsg interface{}) {
	p.Name = Name
	p.UnionId = UnionId
	p.ExtraMsg = ExtraMsg
}

// 获取玩家客户端Id
func (p *Player) ClientId() int32 {
	return p.clientId
}

// 设置玩家客户端Id
func (p *Player) SetClientId(clientId int32) {
	p.clientId = clientId
}
