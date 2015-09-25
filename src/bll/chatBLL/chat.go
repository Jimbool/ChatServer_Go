package chatBLL

import (
	"encoding/json"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/sensitiveWordsBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/webBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/channelType"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/ChatServer_Go/src/model/commandType"
	"github.com/Jordanzuo/ChatServer_Go/src/model/player"
	"github.com/Jordanzuo/ChatServer_Go/src/model/responseDataObject"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/securityUtil"
	"github.com/Jordanzuo/goutil/stringUtil"
	"net"
	"time"
)

var (
	// 客户端连接列表
	ClientList = make(map[*net.Conn]*client.Client)

	// 玩家列表
	PlayerList = make(map[string]*player.Player)

	// 定义增加、删除客户端channel；增加、删除玩家的channel
	ClientAddChan    = make(chan *player.PlayerAndClient)
	ClientRemoveChan = make(chan *player.PlayerAndClient)
	PlayerAddChan    = make(chan *player.PlayerAndClient, 50)
	PlayerRemoveChan = make(chan *player.PlayerAndClient, 50)
)

func init() {
	// 启动处理增加、删除客户端channel；增加、删除玩家的channel的gorountine
	go handleChannel(ClientAddChan, ClientRemoveChan, PlayerAddChan, PlayerRemoveChan)

	// 启动清理过期客户端连接的gorountine
	go clearExpiredClient(ClientRemoveChan)
}

// 处理增加、删除客户端channel；增加、删除玩家的channel的逻辑
// clientAddChan: 客户端增加的channel
// clientRemoveChan: 客户端移除的channel
// playerAddChan: 玩家增加的channel
// playerRemoveChan: 玩家移除的channel
func handleChannel(clientAddChan, clientRemoveChan, playerAddChan, playerRemoveChan chan *player.PlayerAndClient) {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
		}
	}()

	for {
		select {
		case playerAndClientObj := <-clientAddChan:
			addClient(playerAndClientObj.Client)
		case playerAndClientObj := <-clientRemoveChan:
			removeClient(playerAndClientObj.Client)
		case playerAndClientObj := <-playerAddChan:
			addPlayer(playerAndClientObj.Client, playerAndClientObj.Player)
		case playerAndClientObj := <-playerRemoveChan:
			removePlayer(playerAndClientObj.Client, playerAndClientObj.Player)
		case message := <-webBLL.MessageCh:
			broadcast(message)
		default:
			// 休眠一下，防止CPU过高
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// 清理过期的客户端
func clearExpiredClient(clientRemoveChan chan *player.PlayerAndClient) {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
		}
	}()

	for {
		// 休眠指定的时间（单位：秒）(放在此处是因为程序刚启动时并没有过期的客户端，所以先不用占用资源；并且此时LogPath尚未设置，如果直接执行后面的代码会出现panic异常)
		time.Sleep(configBLL.CheckExpiredInterval * time.Second)

		// 清理之前的客户端数量和玩家数量
		beforeClientCount := len(ClientList)
		beforePlayerCount := len(PlayerList)

		// 获取过期的客户端列表
		expiredClientList := getExpiredClientList()

		// 获取本次清理的客户端数量
		expiredClientCount := len(expiredClientList)
		if expiredClientCount == 0 {
			continue
		}

		// 记录日志
		logUtil.Log(fmt.Sprintf("清理前的客户端数量为：%d， 清理前的玩家数量为：%d， 本次清理不活跃的数量为：%d", beforeClientCount, beforePlayerCount, expiredClientCount), logUtil.Debug, true)

		// 移除过期的客户
		for _, item := range expiredClientList {
			item.Quit()
		}
	}
}

// 获取过期的客户端对象列表
// 返回值：过期的客户端对象列表
func getExpiredClientList() (expiredClientList []*client.Client) {
	for _, item := range ClientList {
		if item.HasExpired() {
			expiredClientList = append(expiredClientList, item)
		}
	}

	return
}

// 添加一个新的客户端对象到列表中
// clientObj：客户端对象
func addClient(clientObj *client.Client) {
	ClientList[clientObj.Id] = clientObj
}

// 移除一个客户端对象
// 由于客户端对象与玩家对象之间可能已经建立了对应关系，所以在移除完客户端对象后，还需要移除客户端对象和玩家对象之间的对应关系（如果存在对应关系）
// clientObj：客户端对象
func removeClient(clientObj *client.Client) {
	// 清除在PlayerList中的玩家对象
	if clientObj.PlayerId != "" {
		delete(PlayerList, clientObj.PlayerId)
	}

	// 删除客户端
	delete(ClientList, clientObj.Id)
}

// 添加玩家对象
// 添加玩家的时候，客户端对象已经存在，所以在添加完玩家对象后，还需要添加客户端对象与玩家对象之间的对应关系
// clientObj：客户端对象
// playerObj：玩家对象
func addPlayer(clientObj *client.Client, playerObj *player.Player) {
	// 如果玩家Id有对应的玩家对象(也就是所谓的重复登陆或非正常途径退出)
	if oldPlayerObj, ok := PlayerList[playerObj.Id]; ok {
		if oldClientObj, ok := ClientList[oldPlayerObj.ClientId]; ok {
			// 必须要先Logout再Quit，否则新登陆的玩家也会显示未登陆
			oldClientObj.PlayerLogout()
			oldClientObj.Quit()
		}
	}

	// 更新Client对象对应的PlayerId
	clientObj.PlayerLogin(playerObj.Id)

	// 将玩家对象添加到PlayerList列表中
	PlayerList[playerObj.Id] = playerObj
}

// 移除客户端对象
// removeClient和removePlayer的区别在于：
// removeClient的在移除Client对象、以及Client与Player的对应关系，同时要将对应的Player也移除，但
// removePlayer只移除Player以及Player与Client的对应关系，Client对象会继续保留
// clientObj：客户端对象
// playerObj：玩家对象
func removePlayer(clientObj *client.Client, playerObj *player.Player) {
	// 清除在PlayerList中的玩家对象
	delete(PlayerList, playerObj.Id)

	// 将对应的Client对象的玩家Id清空
	clientObj.PlayerLogout()
}

// 根据客户端对象获取对应的玩家对象
// clientObj：客户端对象
// 返回值：玩家对象
func getPlayerByClient(clientObj *client.Client) (*player.Player, bool) {
	// 先根据客户端Id从ClientAndPlayerList找到对应的PlayerId
	if clientObj.PlayerId != "" {
		// 然后根据PlayerId从PlayerList找到对应的玩家对象
		if playerObj, ok := PlayerList[clientObj.PlayerId]; ok {
			return playerObj, true
		}
	}

	return nil, false
}

// 根据玩家对象获取对应的客户端对象
// playerObj：玩家对象
// 返回值：客户端对象
func getClientByPlayer(playerObj *player.Player) (*client.Client, bool) {
	// 先根据PlayerId从PlayerAndClientList获得ClientId
	if playerObj.ClientId != nil {
		// 再根据ClientId从ClientList中获得对应的Client对象
		if clientObj, ok := ClientList[playerObj.ClientId]; ok {
			return clientObj, true
		}
	}

	return nil, false
}

// 获取初始的响应对象
// ct：请求的命令类型
// 返回值：响应对象
func getInitResponseObj(ct commandType.CommandType) responseDataObject.ResponseObject {
	return responseDataObject.ResponseObject{
		Code:        responseDataObject.Success,
		Message:     "",
		Data:        nil,
		CommandType: ct,
	}
}

// 获取响应类型为数据错误的响应对象
// responseObj：响应对象
// 返回值：响应对象
func getDataErrorReponseObj(responseObj responseDataObject.ResponseObject) responseDataObject.ResponseObject {
	return getResultStatusResponseObj(responseObj, responseDataObject.DataError)
}

// 获取指定响应类型的响应对象
// responseObj：响应对象
// rs：响应类型对象
// 返回值：响应对象
func getResultStatusResponseObj(responseObj responseDataObject.ResponseObject, rs responseDataObject.ResultStatus) responseDataObject.ResponseObject {
	responseObj.Code = rs
	responseObj.Message = rs.String()

	return responseObj
}

// 处理请求
// clientObj：对应的客户端对象
// request：请求内容字节数组(json格式)
// 返回值：无
func HanleRequest(clientObj *client.Client, request []byte, clientAddChan, clientRemoveChan, playerAddChan, playerRemoveChan chan *player.PlayerAndClient) {
	responseObj := getInitResponseObj(commandType.Login)

	// 最后将responseObject发送到客户端
	defer func() {
		// 如果不成功，则向客户端发送数据；如果成功，则已经通过对应的方法发送结果，故不通过此处
		if responseObj.Code != responseDataObject.Success {
			responseResult(clientObj, responseObj)
		}
	}()

	// 解析请求字符串
	requestMap := make(map[string]interface{})
	err := json.Unmarshal(request, &requestMap)
	if err != nil {
		logUtil.Log(fmt.Sprintf("反序列化%s出错，错误信息为：%s", string(request), err), logUtil.Error, true)
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	// 解析CommandType
	var ok bool
	commandType_int, ok := requestMap["CommandType"].(float64)
	if !ok {
		logUtil.Log(fmt.Sprintf("CommandType:%v，不是int类型", requestMap["CommandType"]), logUtil.Error, true)
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	// 得到真实的CommandType
	commandType_real := commandType.CommandType(int(commandType_int))

	// 设置responseObject的CommandType
	responseObj.CommandType = commandType_real

	// 定义Player对象
	var playerObj *player.Player

	// 如果不是Login方法，则判断Client对象所对应的玩家对象是否存在（因为当是Login方法时，Player对象尚不存在）
	if commandType_real != commandType.Login {
		if playerObj, ok = getPlayerByClient(clientObj); !ok {
			responseObj = getResultStatusResponseObj(responseObj, responseDataObject.NoLogin)
			return
		}
	}

	// 解析Command(是map[string]interface{}类型)
	commandMap, ok := requestMap["Command"].(map[string]interface{})
	if !ok {
		logUtil.Log(fmt.Sprintf("commandMap:%v，不是map类型", commandMap), logUtil.Error, true)
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	// 根据不同的请求方法，来调用不同的处理方式
	switch commandType_real {
	case commandType.Login:
		responseObj = login(clientObj, commandType_real, commandMap, playerAddChan)
	case commandType.Logout:
		responseObj = logout(clientObj, playerObj, commandType_real, playerRemoveChan)
	case commandType.SendMessage:
		responseObj = sendMessage(clientObj, playerObj, commandType_real, commandMap)
	case commandType.UpdatePlayerInfo:
		responseObj = updatePlayerInfo(clientObj, playerObj, commandType_real, commandMap)
	default:
		responseObj = getResultStatusResponseObj(responseObj, responseDataObject.CommandTypeNotDefined)
	}
}

func login(clientObj *client.Client, ct commandType.CommandType, commandMap map[string]interface{}, playerAddChan chan *player.PlayerAndClient) (responseObj responseDataObject.ResponseObject) {
	responseObj = getInitResponseObj(ct)

	// 解析参数
	var ok bool
	var id string
	var name string
	var unionId string
	var sign string
	var extraMsg interface{}

	id, ok = commandMap["Id"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("Id:%v，不是string类型", commandMap["Id"]), logUtil.Error, true)
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	name, ok = commandMap["Name"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("Name:%v，不是string类型", commandMap["Name"]), logUtil.Error, true)
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	unionId, ok = commandMap["UnionId"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("UnionId:%v，不是string类型", commandMap["UnionId"]), logUtil.Error, true)
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	sign, ok = commandMap["Sign"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("Sign:%v，不是string类型", commandMap["Sign"]), logUtil.Error, true)
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	extraMsg = commandMap["ExtraMsg"]

	// 验证签名是否正确
	rawstring := fmt.Sprintf("%s-%s-%s", id, name, configBLL.LoginKey)
	if sign != securityUtil.Md5String(rawstring, false) {
		responseObj = getResultStatusResponseObj(responseObj, responseDataObject.SignError)
		return
	}

	// 构造玩家对象
	playerObj := player.NewPlayer(id, name, unionId, extraMsg, clientObj.Id)

	// 将玩家对象添加到玩家增加的channel中
	playerAddChan <- player.NewPlayerAndClient(playerObj, clientObj)

	// 输出结果
	responseResult(clientObj, responseObj)

	return
}

func logout(clientObj *client.Client, playerObj *player.Player, ct commandType.CommandType, playerRemoveChan chan *player.PlayerAndClient) (responseObj responseDataObject.ResponseObject) {
	responseObj = getInitResponseObj(ct)

	// 将玩家对象添加到玩家移除的channel中
	playerRemoveChan <- player.NewPlayerAndClient(playerObj, clientObj)

	// 输出结果
	responseResult(clientObj, responseObj)

	return
}

func updatePlayerInfo(clientObj *client.Client, playerObj *player.Player, ct commandType.CommandType, commandMap map[string]interface{}) (responseObj responseDataObject.ResponseObject) {
	responseObj = getInitResponseObj(ct)

	// 解析参数
	var ok bool
	var name string
	var unionId string
	var extraMsg interface{}

	name, ok = commandMap["Name"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("Name:%v，不是string类型", commandMap["Name"]), logUtil.Error, true)
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	unionId, ok = commandMap["UnionId"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("UnionId:%v，不是string类型", commandMap["UnionId"]), logUtil.Error, true)
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	extraMsg = commandMap["ExtraMsg"]

	// 更新玩家信息
	playerObj.Update(name, unionId, extraMsg)

	// 输出结果
	responseResult(clientObj, responseObj)

	return
}

func sendMessage(clientObj *client.Client, playerObj *player.Player, ct commandType.CommandType, commandMap map[string]interface{}) (responseObj responseDataObject.ResponseObject) {
	responseObj = getInitResponseObj(ct)

	// 解析参数
	var ok bool
	var channelType_real channelType.ChannelType
	var message string

	channelType_int, ok := commandMap["ChannelType"].(float64)
	if !ok {
		logUtil.Log(fmt.Sprintf("ChannelType:%v，不是int类型", commandMap["ChannelType"]), logUtil.Error, true)
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	// 得到真实的ChannelType
	channelType_real = channelType.ChannelType(int(channelType_int))

	message, ok = commandMap["Message"].(string)
	if !ok {
		logUtil.Log(fmt.Sprintf("Message:%v，不是string类型", commandMap["Message"]), logUtil.Error, true)
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	// 判断消息长度是否超过最大值，如果超过最大值，则只趣前面部分
	if len(message) > configBLL.MaxMsgLength {
		message = stringUtil.Substring(message, 0, configBLL.MaxMsgLength)
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
		for _, item := range PlayerList {
			finalPlayerList = append(finalPlayerList, item)
		}
	case channelType.Union:
		// 判断公会Id是否为空
		if playerObj.UnionId == "" {
			responseObj = getResultStatusResponseObj(responseObj, responseDataObject.NotInUnion)
			return
		}

		// 筛选同一个公会的成员
		for _, item := range PlayerList {
			if playerObj.UnionId == item.UnionId {
				finalPlayerList = append(finalPlayerList, item)
			}
		}
	case channelType.Private:
		toPlayerId, ok := commandMap["ToPlayerId"].(string)
		if !ok {
			logUtil.Log(fmt.Sprintf("ToPlayerId:%v，不是string类型", commandMap["ToPlayerId"]), logUtil.Error, true)
			responseObj = getDataErrorReponseObj(responseObj)
			return
		}

		// 不能给自己发送消息
		if playerObj.Id == toPlayerId {
			responseObj = getResultStatusResponseObj(responseObj, responseDataObject.CantSendMessageToSelf)
			return
		}

		// 获得目标玩家对象
		toPlayerObj, ifToPlayerExists = PlayerList[toPlayerId]
		if !ifToPlayerExists {
			responseObj = getResultStatusResponseObj(responseObj, responseDataObject.NotFoundTarget)
			return
		}

		// 添加到列表中
		finalPlayerList = append(finalPlayerList, playerObj, toPlayerObj)
	default:
		responseObj = getDataErrorReponseObj(responseObj)
		return
	}

	// 组装需要发送的数据
	data := make(map[string]interface{})
	data["ChannelType"] = channelType_real
	data["Message"] = message

	// 增加发送者信息
	from := make(map[string]interface{})
	from["Id"] = playerObj.Id
	from["Name"] = playerObj.Name
	from["UnionId"] = playerObj.UnionId
	from["ExtraMsg"] = playerObj.ExtraMsg

	data["From"] = from

	// 如果是私聊，则加上私聊对象的信息
	if ifToPlayerExists {
		to := make(map[string]interface{})
		to["Id"] = toPlayerObj.Id
		to["Name"] = toPlayerObj.Name
		to["UnionId"] = toPlayerObj.UnionId
		to["ExtraMsg"] = toPlayerObj.ExtraMsg

		data["To"] = to
	}

	// 设置responseObj的Data属性
	responseObj.Data = data

	// 遍历，向玩家发送消息
	for _, item := range finalPlayerList {
		// 根据Player对象获得Client对象
		if finalClientObj, ok := getClientByPlayer(item); ok {
			responseResult(finalClientObj, responseObj)
		}
	}

	return
}

// 广播消息
func broadcast(message string) {
	fmt.Println("broadcast:", message)

	responseObj := getInitResponseObj(commandType.SendMessage)

	// 组装需要发送的数据
	data := make(map[string]interface{})
	data["ChannelType"] = channelType.World
	data["Message"] = message

	// 增加发送者信息
	from := make(map[string]interface{})
	from["Id"] = "system"
	from["Name"] = "系统"
	from["UnionId"] = ""
	from["ExtraMsg"] = ""

	data["From"] = from

	responseObj.Data = data

	// 遍历，向玩家发送消息
	for _, item := range PlayerList {
		// 根据Player对象获得Client对象
		if finalClientObj, ok := getClientByPlayer(item); ok {
			responseResult(finalClientObj, responseObj)
		}
	}
}

// 发送响应结果
// clientObj：客户端对象
// responseObject：响应对象
func responseResult(clientObj *client.Client, responseObject responseDataObject.ResponseObject) {
	b, err := json.Marshal(responseObject)
	if err != nil {
		logUtil.Log(fmt.Sprintf("序列化输出结果%v出错", responseObject), logUtil.Error, true)
	} else {
		clientObj.SendByteMessage(b)
	}
}
