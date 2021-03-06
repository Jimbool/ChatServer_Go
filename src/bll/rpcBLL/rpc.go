/*
远程服务的包，用于提供Sokcet服务
*/
package rpcBLL

import (
	"errors"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/chatBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/clientBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/playerBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/ChatServer_Go/src/model/disconnectType"
	"github.com/Jordanzuo/goutil/logUtil"
	"net"
	"sync"
)

// 处理客户端逻辑
// clientObj：客户端对象
func handleClientContent(clientObj *client.Client) {
	for {
		content, ok := clientObj.GetValidMessage()
		if !ok {
			break
		}

		// 处理数据，如果长度为0则表示心跳包
		if len(content) == 0 {
			continue
		} else {
			chatBLL.HanleRequest(clientObj, content)
		}
	}
}

// 处理客户端连接
// conn：客户端连接对象
func handleConn(conn net.Conn) {
	// 创建客户端对象
	clientObj := client.NewClient(conn)

	// 将客户端对象添加到客户端增加的channel中
	clientBLL.RegisterClient(clientObj)

	// 将客户端对象添加到客户端移除的channel中
	defer func() {
		playerBLL.DisconnectByClient(clientObj, disconnectType.FromRpc)
	}()

	// 无限循环，不断地读取数据，解析数据，处理数据
	for {
		// 先读取数据，每次读取1024个字节
		readBytes := make([]byte, 1024)

		// Read方法会阻塞，所以不用考虑异步的方式
		n, err := conn.Read(readBytes)
		if err != nil {
			break
		}

		// 将读取到的数据追加到已获得的数据的末尾
		clientObj.AppendContent(readBytes[:n])

		// 处理数据
		handleClientContent(clientObj)
	}
}

// 启动服务器
func StartServer(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	logUtil.Log("Socket服务器开始监听...", logUtil.Info, true)

	// 监听指定的端口
	listener, err := net.Listen("tcp", configBLL.SocketServerAddress())
	if err != nil {
		panic(errors.New(fmt.Sprintf("Listen Error: %s", err)))
	} else {
		msg := fmt.Sprintf("Got listener for the server. (local address: %s)", listener.Addr())

		// 记录和显示日志，并且判断是否需要退出
		logUtil.Log(msg, logUtil.Info, true)
		fmt.Println(msg)
	}

	for {
		// 阻塞直至新连接到来
		conn, err := listener.Accept()
		if err != nil {
			logUtil.Log(fmt.Sprintf("Accept Error: %s", err), logUtil.Error, true)
			continue
		}

		// 启动一个新协程来处理链接
		go handleConn(conn)
	}
}
