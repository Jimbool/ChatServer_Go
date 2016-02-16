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
	"sync"
	"time"
)

var (
	historyMessageList = make([]*responseDataObject.SocketResponseObject, 0, 20)
	mutex              sync.RWMutex
)

func addNewMessage(responseObj *responseDataObject.SocketResponseObject) {
	mutex.Lock()
	defer mutex.Unlock()

	historyMessageList = append(historyMessageList, responseObj)
	if len(historyMessageList) > configBLL.MaxHistoryCount() {
		historyMessageList = historyMessageList[len(historyMessageList)-configBLL.MaxHistoryCount():]
	}
}

func pushMessageAfterLogin(clientObj *client.Client) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, message := range historyMessageList {
		playerBLL.SendToClient(clientObj, message)
	}
}

// 处理客户端请求
// clientObj：对应的客户端对象
// request：请求内容字节数组(json格式)
// 返回值：无
func HanleRequest(clientObj *client.Client, request []byte) {
	start := time.Now().Unix()
	responseObj := responseDataObject.NewSocketResponseObject(commandType.Login)

	defer func() {
		// 如果不成功，则向客户端发送数据；因为成功已经通过对应的方法发送结果，故不通过此处
		if responseObj.Code != responseDataObject.Success {
			// 如果是客户端数据错误，则将客户端请求数据记录下来
			if responseObj.Code == responseDataObject.ClientDataError {
				logUtil.Log(fmt.Sprintf("请求的数据为：%s, 返回的结果为客户端数据错误", string(request)), logUtil.Error, true)
			}

			playerBLL.SendToClient(clientObj, responseObj)
		}

		// 如果处理的时间超过3s，则记录下来以便于后续分析
		end := time.Now().Unix()
		duration := end - start
		if duration > 3 {
			logUtil.Log(fmt.Sprintf("请求内容为：%s，请求时间为%d秒", string(request), duration), logUtil.Warn, true)
		}
	}()

	// 定义变量
	var requestMap map[string]interface{}
	var commandMap map[string]interface{}
	var playerObj *player.Player
	var exists bool
	var ok bool
	var err error

	// 解析请求字符串
	if err = json.Unmarshal(request, &requestMap); err != nil {
		logUtil.Log(fmt.Sprintf("反序列化出错，错误信息为：%s", err), logUtil.Error, true)
		responseObj.SetClientDataError()
		return
	}

	// 解析CommandType
	if commandType_float, ok := requestMap["CommandType"].(float64); !ok {
		logUtil.Log("CommandType不是int类型", logUtil.Error, true)
		responseObj.SetClientDataError()
		return
	} else {
		// 设置responseObject的CommandType
		responseObj.SetCommandType(commandType.CommandType(int(commandType_float)))
	}

	// 如果不是Login方法，则判断Client对象所对应的玩家对象是否存在（因为当是Login方法时，Player对象尚不存在）
	if responseObj.CommandType != commandType.Login {
		playerObj, exists, err = playerBLL.GetPlayer(clientObj.PlayerId(), false)
		if err != nil {
			responseObj.SetDataError()
			return
		}

		if !exists {
			responseObj.SetResultStatus(responseDataObject.NoLogin)
			return
		}
	}

	// 解析Command(是map[string]interface{}类型)；只有当不是Logout方法时才解析，因为Logout时Command为空
	if responseObj.CommandType != commandType.Logout {
		if commandMap, ok = requestMap["Command"].(map[string]interface{}); !ok {
			logUtil.Log("commandMap不是map类型", logUtil.Error, true)
			responseObj.SetClientDataError()
			return
		}
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

	// 定义变量
	var ok bool
	var exists bool
	var id string
	var name string
	var unionId string
	var sign string
	var extraMsg string
	var isNewPlayer bool
	var err error
	var playerObj *player.Player
	var gamePlayerName string
	var gameUnionId string

	if id, ok = commandMap["Id"].(string); !ok {
		logUtil.Log("Id不是string类型", logUtil.Error, true)
		responseObj.SetClientDataError()
		return responseObj
	}

	if name, ok = commandMap["Name"].(string); !ok {
		logUtil.Log("Name不是string类型", logUtil.Error, true)
		responseObj.SetClientDataError()
		return responseObj
	}

	if unionId, ok = commandMap["UnionId"].(string); !ok {
		logUtil.Log("UnionId不是string类型", logUtil.Error, true)
		responseObj.SetClientDataError()
		return responseObj
	}

	if sign, ok = commandMap["Sign"].(string); !ok {
		logUtil.Log("Sign不是string类型", logUtil.Error, true)
		responseObj.SetClientDataError()
		return responseObj
	}

	if extraMsg, ok = commandMap["ExtraMsg"].(string); !ok {
		logUtil.Log("ExtraMsg不是string类型", logUtil.Error, true)
		responseObj.SetClientDataError()
		return responseObj
	}

	// 验证签名是否正确
	if verifySign(id, name, sign) == false {
		responseObj.SetResultStatus(responseDataObject.SignError)
		return responseObj
	}

	// 判断玩家是否在缓存中已经存在
	playerObj, exists, err = playerBLL.GetPlayer(id, false)
	if err != nil {
		responseObj.SetDataError()
		return responseObj
	}

	if exists {
		name = playerObj.Name

		// 判断是否重复登陆
		if playerObj.ClientId > 0 {
			if oldClientObj, exists := clientBLL.GetClient(playerObj.ClientId); exists {
				// 如果不是同一个客户端，则先给客户端发送在其他设备登陆信息，然后断开连接
				if clientObj != oldClientObj {
					playerBLL.SendLoginAnotherDeviceMsg(oldClientObj)
				}
			}
		}
	} else {
		// 判断数据库中是否已经存在该玩家，如果不存在则表明是新玩家，先到游戏库中验证
		playerObj, exists, err = playerBLL.GetPlayer(id, true)
		if err != nil {
			responseObj.SetDataError()
			return responseObj
		}

		if !exists {
			// 验证玩家Id在游戏库中是否存在
			gamePlayerName, gameUnionId, exists, err = playerBLL.GetGamePlayer(id)
			if err != nil {
				responseObj.SetDataError()
				return responseObj
			} else if !exists {
				responseObj.SetResultStatus(responseDataObject.PlayerNotExist)
				return responseObj
			} else {
				if name != gamePlayerName {
					responseObj.SetResultStatus(responseDataObject.NameError)
					return responseObj
				}

				if unionId != "" && unionId != "00000000-0000-0000-0000-000000000000" && unionId != gameUnionId {
					responseObj.SetResultStatus(responseDataObject.UnionIdError)
					return responseObj
				}
			}

			if playerObj, err = playerBLL.RegisterNewPlayer(id, name, unionId, extraMsg); err != nil {
				responseObj.SetDataError()
				return responseObj
			}
			isNewPlayer = true
		}
	}

	// 判断玩家是否被封号
	if playerObj.IsForbidden {
		responseObj.SetResultStatus(responseDataObject.PlayerIsForbidden)
		return responseObj
	}

	// 更新客户端对象的玩家Id
	clientObj.PlayerLogin(id)

	// 更新玩家登录信息
	if err = playerBLL.UpdateLoginInfo(playerObj, clientObj, isNewPlayer); err != nil {
		responseObj.SetDataError()
		return responseObj
	}

	// 将玩家对象添加到玩家列表中
	playerBLL.RegisterPlayer(playerObj)

	// 输出结果
	playerBLL.SendToClient(clientObj, responseObj)

	// 推送历史信息
	pushMessageAfterLogin(clientObj)

	return responseObj
}

func logout(clientObj *client.Client, playerObj *player.Player, ct commandType.CommandType) *responseDataObject.SocketResponseObject {
	responseObj := responseDataObject.NewSocketResponseObject(ct)

	// 玩家登出
	clientObj.LogoutAndQuit()

	// 将玩家对象从缓存中移除
	playerBLL.UnRegisterPlayer(playerObj)

	return responseObj
}

func updatePlayerInfo(clientObj *client.Client, playerObj *player.Player, ct commandType.CommandType, commandMap map[string]interface{}) *responseDataObject.SocketResponseObject {
	responseObj := responseDataObject.NewSocketResponseObject(ct)

	// 定义变量
	var exists bool
	var ok bool
	var name string
	var unionId string
	var extraMsg string
	var err error
	var gamePlayerName string
	var gameUnionId string

	if name, ok = commandMap["Name"].(string); !ok {
		logUtil.Log("Name不是string类型", logUtil.Error, true)
		responseObj.SetClientDataError()
		return responseObj
	}

	if unionId, ok = commandMap["UnionId"].(string); !ok {
		logUtil.Log("UnionId不是string类型", logUtil.Error, true)
		responseObj.SetClientDataError()
		return responseObj
	}

	if extraMsg, ok = commandMap["ExtraMsg"].(string); !ok {
		logUtil.Log("ExtraMsg不是string类型", logUtil.Error, true)
		responseObj.SetClientDataError()
		return responseObj
	}

	// 如果玩家名或公会Id有改变，则到游戏库中去验证是否是正确的名称
	if name != playerObj.Name || unionId != playerObj.UnionId {
		// 验证玩家Id在游戏库中是否存在
		gamePlayerName, gameUnionId, exists, err = playerBLL.GetGamePlayer(playerObj.Id)
		if err != nil {
			responseObj.SetDataError()
			return responseObj
		} else if !exists {
			responseObj.SetResultStatus(responseDataObject.PlayerNotExist)
			return responseObj
		} else {
			if name != gamePlayerName {
				responseObj.SetResultStatus(responseDataObject.NameError)
				return responseObj
			}

			if unionId != "" && unionId != "00000000-0000-0000-0000-000000000000" && unionId != gameUnionId {
				responseObj.SetResultStatus(responseDataObject.UnionIdError)
				return responseObj
			}
		}
	}

	// 更新玩家信息
	if err = playerBLL.UpdateInfo(playerObj, name, unionId, extraMsg); err != nil {
		responseObj.SetDataError()
		return responseObj
	}

	// 输出结果
	playerBLL.SendToClient(clientObj, responseObj)

	return responseObj
}

func sendMessage(clientObj *client.Client, playerObj *player.Player, ct commandType.CommandType, commandMap map[string]interface{}) *responseDataObject.SocketResponseObject {
	responseObj := responseDataObject.NewSocketResponseObject(ct)

	// 判断玩家是否被禁言
	if isInSilent, _ := playerObj.IsInSilent(); isInSilent {
		responseObj.SetResultStatus(responseDataObject.PlayerIsInSilent)
		return responseObj
	}

	// 定义变量
	var ok bool
	var channelType_real channelType.ChannelType
	var message string
	var err error
	var finalPlayerList = make([]*player.Player, 0, 1024)
	var toPlayerId string
	var toPlayerObj *player.Player
	var ifToPlayerExists bool

	if channelType_float, ok := commandMap["ChannelType"].(float64); !ok {
		logUtil.Log("ChannelType不是int类型", logUtil.Error, true)
		responseObj.SetClientDataError()
		return responseObj
	} else {
		channelType_real = channelType.ChannelType(int(channelType_float))
	}

	if message, ok = commandMap["Message"].(string); !ok {
		logUtil.Log("Message不是string类型", logUtil.Error, true)
		responseObj.SetClientDataError()
		return responseObj
	}

	// 判断消息长度是否超过最大值，如果超过最大值，则只趣前面部分
	if len(message) > configBLL.MaxMsgLength() {
		message = stringUtil.Substring(message, 0, configBLL.MaxMsgLength())
	}

	// 处理敏感词汇
	message = sensitiveWordsBLL.HandleSensitiveWords(message)

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
		if toPlayerId, ok = commandMap["ToPlayerId"].(string); !ok {
			logUtil.Log("ToPlayerId不是string类型", logUtil.Error, true)
			responseObj.SetClientDataError()
			return responseObj
		}

		// 不能给自己发送消息
		if playerObj.Id == toPlayerId {
			responseObj.SetResultStatus(responseDataObject.CantSendMessageToSelf)
			return responseObj
		}

		// 获得目标玩家对象
		toPlayerObj, ifToPlayerExists, err = playerBLL.GetPlayer(toPlayerId, false)
		if err != nil {
			responseObj.SetResultStatus(responseDataObject.DataError)
			return responseObj
		}

		if !ifToPlayerExists {
			responseObj.SetResultStatus(responseDataObject.NotFoundTarget)
			return responseObj
		}

		// 添加到列表中
		finalPlayerList = append(finalPlayerList, playerObj, toPlayerObj)
	default:
		responseObj.SetClientDataError()
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

	// 如果是世界频道信息，添加到历史消息里面
	if channelType_real == channelType.World {
		addNewMessage(responseObj)
	}

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
