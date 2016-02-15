package responseDataObject

// 服务端响应结果的状态对象，成功是0，非成功以负数来表示
type ResultStatus int

// 定义所有的响应结果的状态枚举值，此种实现方式是GO语言的标准方式
const (
	// 成功
	Success ResultStatus = -1 * iota

	// 数据错误
	DataError

	// API数据错误
	APIDataError

	// 客户端数据错误
	ClientDataError

	// 命令类型未定义
	CommandTypeNotDefined

	// 签名错误
	SignError

	// 尚未登陆
	NoLogin

	// 不在公会中
	NotInUnion

	// 未找到目标
	NotFoundTarget

	// 不能给自己发消息
	CantSendMessageToSelf

	// 玩家不存在
	PlayerNotExist

	// 玩家被封号
	PlayerIsForbidden

	// 玩家被禁言
	PlayerIsInSilent

	// 只支持POST
	OnlySupportPOST

	// API未定义
	APINotDefined

	// 在另一台设备上登录
	LoginOnAnotherDevice

	// 名称错误
	NameError

	// 公会Id错误
	UnionIdError
)

// 定义所有的响应结果的状态值所对应的字符串描述信息，如果要增加状态枚举，则此处也要相应地增加
var status = [...]string{
	"Success",
	"DataError",
	"APIDataError",
	"ClientDataError",
	"CommandTypeNotDefined",
	"SignError",
	"NoLogin",
	"NotInUnion",
	"NotFoundTarget",
	"CantSendMessageToSelf",
	"PlayerNotExist",
	"PlayerIsForbidden",
	"PlayerIsInSilent",
	"OnlySupportPOST",
	"APINotDefined",
	"LoginOnAnotherDevice",
	"NameError",
	"UnionIdError",
}

// 返回响应状态枚举值对应的描述信息字符串
func (rs ResultStatus) String() string {
	return status[rs*-1]
}
