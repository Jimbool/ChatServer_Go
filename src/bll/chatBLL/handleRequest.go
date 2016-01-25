package chatBLL

import (
	"encoding/json"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/clientBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/playerBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/sensitiveWordsBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/channelType"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/ChatServer_Go/src/model/commandType"
	"github.com/Jordanzuo/ChatServer_Go/src/model/player"
	"github.com/Jordanzuo/ChatServer_Go/src/model/responseDataObject"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/securityUtil"
	"github.com/Jordanzuo/goutil/stringUtil"
	"time"
)

// 处理客户端请求
// clientObj：对应的客户端对象
// request：请求内容字节数组(json格式)
// 返回值：无
func HanleRequest(clientObj *client.Client, request []byte) {
	responseObj := responseDataObject.NewSocketResponseObject(commandType.Login)

	// 最后将responseObject发送到客户端
	defer func() {
		// 如果不成功，则向客户端发送数据；如果成功，则已经通过对应的方法发送结果，故不通过此处
		if responseObj.Code != responseDataObject.Success {
			playerBLL.SendToClient(clientObj, responseObj)
		}
	}()

	// 解析请求字符串
	requestMap := make(map[string]interface{})
	err := json.Unmarshal(request, &requestMap)
	if err != nil {
		logUtil.Log(fmt.Sprintf("反序列化%s出错，错误信息为：%s", string(request), err), logUtil.Error, true)
		responseObj.SetDataError()
		return
	}

	// 解析CommandType
	var ok bool
	commandType_float, ok := requestMap["CommandType"].(float64)
	if !ok {
		logUtil.Log(fmt.Sprintf("CommandType:%v，不是int类型", requestMap["CommandType"]), logUtil.Error, true)
		responseObj.SetDataError()
		return
	}

	// 设置responseObject的CommandType
	responseObj.SetCommandType(commandType.CommandType(int(commandType_float)))

	// 定义Player对象
	var playerObj *player.Player

	// 如果不是Login方法，则判断Client对象所对应的玩家对象是否存在（因为当是Login方法时，Player对象尚不存在）
	if responseObj.CommandType != commandType.Login {
		if playerObj, ok = playerBLL.GetPlayer(clientObj.PlayerId(), false); !ok {
			responseObj.SetResultStatus(responseDataObject.NoLogin)
			return
		}
	}

	// 解析Command(是map[string]interface{}类型)
	commandMap, ok := requestMap["Command"].(map[string]interface{})
	if !ok {
		logUtil.Log(fmt.Sprintf("commandMap:%v，不是map类型", commandMap), logUtil.Error, true)
		responseObj.SetDataError()
		return
	}

	// 根据不同的请求方法，来调用不同的处理方式
	switch responseObj.CommandType {
	case commandType.Login:
		responseObj = login(clientObj, responseObj.CommandType, commandMap)
	case commandType.Logout:
		responseObj = logout(clientObj, playerObj, responseObj.CommandType)
	case commandType.SendMessage:
		responseObj = sendMessage(clientObj, playerObj, responseObj.CommandType, commandMap)
	case commandType.UpdatePlayerInfo:
		responseObj = updatePlayerInfo(clientObj, playerObj, responseObj.CommandType, commandMap)
	default:
		responseObj.SetResultStatus(responseDataObject.CommandTypeNotDefined)
	}
}

func login(clientObj *client.Client, ct commandType.CommandType, commandMap map[string]interface{}) *responseDataObject.SocketResponseObject {
	responseObj := responseDataObject.NewSocketResponseObject(ct)

	// 解析参数
	var ok bool
	var id string
	var name string
	var unionId string
	var sign string
	var extraMsg string

	id, ok = commandMap["Id"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("Id:%v，不是string类型", commandMap["Id"]), logUtil.Error, true)
		responseObj.SetDataError()
		return responseObj
	}

	name, ok = commandMap["Name"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("Name:%v，不是string类型", commandMap["Name"]), logUtil.Error, true)
		responseObj.SetDataError()
		return responseObj
	}

	unionId, ok = commandMap["UnionId"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("UnionId:%v，不是string类型", commandMap["UnionId"]), logUtil.Error, true)
		responseObj.SetDataError()
		return responseObj
	}

	sign, ok = commandMap["Sign"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("Sign:%v，不是string类型", commandMap["Sign"]), logUtil.Error, true)
		responseObj.SetDataError()
		return responseObj
	}

	extraMsg, ok = commandMap["ExtraMsg"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("ExtraMsg:%v，不是string类型", commandMap["ExtraMsg"]), logUtil.Error, true)
		responseObj.SetDataError()
		return responseObj
	}

	// 验证签名是否正确
	if verifySign(id, name, sign) == false {
		responseObj.SetResultStatus(responseDataObject.SignError)
		return responseObj
	}

	// 判断玩家是否在缓存中已经存在
	var playerObj *player.Player
	if playerObj, ok = playerBLL.GetPlayer(id, false); ok {
		// 判断是否重复登陆
		if playerObj.ClientId > 0 {
			if oldClientObj, ok := clientBLL.GetClient(playerObj.ClientId); ok {
				// 如果不是同一个客户端，则先给客户端发送在其他设备登陆信息，然后断开连接
				if clientObj != oldClientObj {
					playerBLL.SendLoginAnotherDeviceMsg(oldClientObj)
					oldClientObj.LogoutAndQuit()
				}
			}
		}
	} else {
		// 判断数据库中是否已经存在该玩家
		if playerObj, ok = playerBLL.GetPlayer(id, true); !ok {
			playerObj = playerBLL.RegisterNewPlayer(id, name, unionId, extraMsg)
		}
	}

	// 判断玩家是否被封号
	if playerObj.IsForbidden {
		responseObj.SetResultStatus(responseDataObject.PlayerIsForbidden)
		return responseObj
	}

	// 判断玩家是否被禁言
	if playerObj.SilentEndTime.Unix() > time.Now().Unix() {
		responseObj.SetResultStatus(responseDataObject.PlayerIsInSilent)
		return responseObj
	}

	// 更新客户端对象的玩家Id
	clientObj.PlayerLogin(id)

	// 更新玩家登录信息
	playerBLL.UpdateLoginInfo(playerObj, clientObj)

	// 将玩家对象添加到玩家增加的channel中
	playerBLL.RegisterPlayer(playerObj)

	// 输出结果
	playerBLL.SendToClient(clientObj, responseObj)

	return responseObj
}

func logout(clientObj *client.Client, playerObj *player.Player, ct commandType.CommandType) *responseDataObject.SocketResponseObject {
	responseObj := responseDataObject.NewSocketResponseObject(ct)

	// 玩家登出
	clientObj.PlayerLogout()

	// 将玩家对象添加到玩家移除的channel中
	playerBLL.UnRegisterPlayer(playerObj)

	// 输出结果
	playerBLL.SendToClient(clientObj, responseObj)

	return responseObj
}

func updatePlayerInfo(clientObj *client.Client, playerObj *player.Player, ct commandType.CommandType, commandMap map[string]interface{}) *responseDataObject.SocketResponseObject {
	responseObj := responseDataObject.NewSocketResponseObject(ct)

	// 解析参数
	var ok bool
	var name string
	var unionId string
	var extraMsg string

	name, ok = commandMap["Name"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("Name:%v，不是string类型", commandMap["Name"]), logUtil.Error, true)
		responseObj.SetDataError()
		return responseObj
	}

	unionId, ok = commandMap["UnionId"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("UnionId:%v，不是string类型", commandMap["UnionId"]), logUtil.Error, true)
		responseObj.SetDataError()
		return responseObj
	}

	extraMsg, ok = commandMap["ExtraMsg"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("ExtraMsg:%v，不是string类型", commandMap["ExtraMsg"]), logUtil.Error, true)
		responseObj.SetDataError()
		return responseObj
	}

	// 更新玩家信息
	playerBLL.UpdateInfo(playerObj, name, unionId, extraMsg)

	// 输出结果
	playerBLL.SendToClient(clientObj, responseObj)

	return responseObj
}

func sendMessage(clientObj *client.Client, playerObj *player.Player, ct commandType.CommandType, commandMap map[string]interface{}) *responseDataObject.SocketResponseObject {
	responseObj := responseDataObject.NewSocketResponseObject(ct)

	// 解析参数
	var ok bool
	var channelType_real channelType.ChannelType
	var message string

	channelType_float, ok := commandMap["ChannelType"].(float64)
	if !ok {
		logUtil.Log(fmt.Sprintf("ChannelType:%v，不是int类型", commandMap["ChannelType"]), logUtil.Error, true)
		responseObj.SetDataError()
		return responseObj
	}

	// 得到真实的ChannelType
	channelType_real = channelType.ChannelType(int(channelType_float))

	message, ok = commandMap["Message"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("Message:%v，不是string类型", commandMap["Message"]), logUtil.Error, true)
		responseObj.SetDataError()
		return responseObj
	}

	// 判断消息长度是否超过最大值，如果超过最大值，则只趣前面部分
	if len(message) > configBLL.MaxMsgLength() {
		message = stringUtil.Substring(message, 0, configBLL.MaxMsgLength())
	}

	// 处理敏感词汇
	message = sensitiveWordsBLL.HandleSensitiveWords(message)

	// 定义变量
	var finalPlayerList = make([]*player.Player, 0, 1024)
	var ifToPlayerExists = false
	var toPlayerObj *player.Player = nil

	// 根据不同的聊天隐疾调用不同的片方法
	switch channelType_real {
	case channelType.World:
		finalPlayerList = playerBLL.GetPlayerList()
	case channelType.Union:
		// 判断公会Id是否为空
		if playerObj.UnionId == "" {
			responseObj.SetResultStatus(responseDataObject.NotInUnion)
			return responseObj
		}

		finalPlayerList = playerBLL.GetPlayerListInSameUnion(playerObj)
	case channelType.Private:
		toPlayerId, ok := commandMap["ToPlayerId"].(string)
		if !ok {
			logUtil.Log(fmt.Sprintf("ToPlayerId:%v，不是string类型", commandMap["ToPlayerId"]), logUtil.Error, true)
			responseObj.SetDataError()
			return responseObj
		}

		// 不能给自己发送消息
		if playerObj.Id == toPlayerId {
			responseObj.SetResultStatus(responseDataObject.CantSendMessageToSelf)
			return responseObj
		}

		// 获得目标玩家对象
		toPlayerObj, ifToPlayerExists = playerBLL.GetPlayer(toPlayerId, false)
		if !ifToPlayerExists {
			responseObj.SetResultStatus(responseDataObject.NotFoundTarget)
			return responseObj
		}

		// 添加到列表中
		finalPlayerList = append(finalPlayerList, playerObj, toPlayerObj)
	default:
		responseObj.SetDataError()
		return responseObj
	}

	// 组装需要发送的数据
	data := make(map[string]interface{})
	data["ChannelType"] = channelType_real
	data["Message"] = message
	data["From"] = playerObj

	// 如果是私聊，则加上私聊对象的信息
	if ifToPlayerExists {
		data["To"] = toPlayerObj
	}

	// 设置responseObj的Data属性
	responseObj.SetData(data)

	// 向玩家发送消息
	playerBLL.SendToPlayer(finalPlayerList, responseObj)

	return responseObj
}

// 验证签名
// id：玩家Id
// name：玩家名称
// sign：签名数据
func verifySign(id, name, sign string) bool {
	rawstring := fmt.Sprintf("%s-%s-%s", id, name, configBLL.AppKey())
	if sign == securityUtil.Md5String(rawstring, false) {
		return true
	}

	return false
}
