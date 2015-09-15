package responseDataObject

import (
	"github.com/Jordanzuo/chatServer/src/model/commandType"
)

// 响应对象
type ResponseObject struct {
	// 响应结果的状态值
	Code ResultStatus

	// 响应结果的状态值所对应的描述信息
	Message string

	// 响应结果的数据
	Data interface{}

	// 响应结果对应的请求命令类型
	CommandType commandType.CommandType
}
