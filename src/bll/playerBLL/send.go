package playerBLL

import (
	"encoding/json"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/clientBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/channelType"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/ChatServer_Go/src/model/commandType"
	"github.com/Jordanzuo/ChatServer_Go/src/model/player"
	"github.com/Jordanzuo/ChatServer_Go/src/model/responseDataObject"
	"github.com/Jordanzuo/goutil/logUtil"
)

// 服务器推送信息
// message：推送的消息
func PushMessage(message string) {
	responseObj := responseDataObject.NewSocketResponseObject(commandType.SendMessage)

	// 组装需要发送的数据
	data := make(map[string]interface{})
	data["ChannelType"] = channelType.World
	data["Message"] = message
	data["From"] = "System"

	responseObj.SetData(data)

	for _, item := range playerList {
		if item.ClientId > 0 {
			if clientObj, ok := clientBLL.GetClient(item.ClientId); ok {
				responseResult(clientObj, responseObj)
			}
		}
	}
}

// 发送在另一台设备登陆的信息
// player：玩家对象
func SendLoginAnotherDevice(clientObj *client.Client) {
	responseObj := responseDataObject.NewSocketResponseObject(commandType.Login)
	responseObj.SetResultStatus(responseDataObject.LoginOnAnotherDevice)
	responseResult(clientObj, responseObj)
}

// 发送数据给客户端
// player：玩家对象
// responseObj：Socket服务器的返回对象
func SendToClient(clientObj *client.Client, responseObj *responseDataObject.SocketResponseObject) {
	responseResult(clientObj, responseObj)
}

// 发送数据给玩家
// playerList：玩家列表
// responseObj：Socket服务器的返回对象
func SendToPlayer(playerList []*player.Player, responseObj *responseDataObject.SocketResponseObject) {
	for _, item := range playerList {
		if item.ClientId > 0 {
			if clientObj, ok := clientBLL.GetClient(item.ClientId); ok {
				responseResult(clientObj, responseObj)
			}
		}
	}
}

// 发送响应结果
// clientObj：客户端对象
// responseObject：响应对象
func responseResult(clientObj *client.Client, responseObj *responseDataObject.SocketResponseObject) {
	b, err := json.Marshal(responseObj)
	if err != nil {
		logUtil.Log(fmt.Sprintf("序列化输出结果%v出错", responseObj), logUtil.Error, true)
	} else {
		clientObj.SendByteMessage(b)
	}
}
