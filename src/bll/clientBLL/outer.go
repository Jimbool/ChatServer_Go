package clientBLL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
)

// 返回过期的客户端列表
// 返回值：
// 过期的客户端列表
func GetExpiredClientList() (expiredClientList []*client.Client) {
	for _, item := range clientList {
		if item.HasExpired() {
			expiredClientList = append(expiredClientList, item)
		}
	}

	return
}

// 添加新的客户端
// clientObj：客户端对象
func RegisterClient(clientObj *client.Client) {
	clientAddChan <- clientObj
}

// 移除客户端
// clientObj：客户端对象
func UnRegisterClient(clientObj *client.Client) {
	clientRemoveChan <- clientObj
}

// 根据客户端Id获取对应的客户端对象
// id：客户端Id
// 返回值：客户端对象
func GetClient(id int32) (*client.Client, bool) {
	if clientObj, ok := clientList[id]; ok {
		return clientObj, true
	}

	return nil, false
}

// 获取客户端的数量
// 返回值：
// 客户端数量
func GetClientCount() int {
	return len(clientList)
}
