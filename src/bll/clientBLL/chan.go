package clientBLL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/goutil/logUtil"
	"time"
)

var (
	// 客户端连接列表
	clientList = make(map[int32]*client.Client)

	// 定义增加、删除客户端channel；增加、删除玩家的channel
	clientAddChan    = make(chan *client.Client)
	clientRemoveChan = make(chan *client.Client)
)

func init() {
	go handleClientChannel()
}

// 处理增加、删除客户端channel；增加、删除玩家的channel的逻辑
func handleClientChannel() {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
		}
	}()

	for {
		select {
		case clientObj := <-clientAddChan:
			addClient(clientObj)
		case clientObj := <-clientRemoveChan:
			removeClient(clientObj)
		default:
			// 休眠一下，防止CPU过高
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// 添加一个新的客户端对象到列表中
// clientObj：客户端对象
func addClient(clientObj *client.Client) {
	clientList[clientObj.Id()] = clientObj
}

// 移除一个客户端对象
// clientObj：客户端对象
func removeClient(clientObj *client.Client) {
	delete(clientList, clientObj.Id())
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
