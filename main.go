package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/chatBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/webBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/client"
	"github.com/Jordanzuo/ChatServer_Go/src/model/player"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/mathUtil"
	"github.com/Jordanzuo/goutil/stringUtil"
	"github.com/Jordanzuo/goutil/timeUtil"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	// 配置文件名称
	CONFIG_FILE_NAME = "config.ini"

	// 服务器网络协议
	SERVER_NETWORK = "tcp"

	// 日志文件路径前缀
	LOG_PATH_SUFFIX = "LOG"
)

var (
	// 服务器监听地址
	ServerAddress string

	// Web服务器监听地址
	WebServerAddress string
)

func init() {
	// 由于服务器的运行依赖于init中执行的逻辑，所以如果出现任何的错误都直接panic，让程序启动失败；而不是让它启动成功，但是在运行时出现错误

	// 读取配置文件（一次性读取整个文件，则使用ioutil）
	bytes, err := ioutil.ReadFile(CONFIG_FILE_NAME)
	if err != nil {
		panic(err)
	}

	// 使用json反序列化
	config := make(map[string]interface{})
	if err = json.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}

	// 本方法内的参数一定需要优先设置，因为这里面的设置是全局的设置，在其它的包内可能会被用到
	// 设置本文件中所需参数
	setParam(config)

	// 设置日志文件的存储目录
	setLogPath()

	// 设置chatBLL的参数
	chatBLL.SetParam(config)

	// 设置client的参数
	client.SetParam(config)
}

// 设置参数
// config：从配置文件里面解析出来的配置内容
func setParam(config map[string]interface{}) {
	// 解析SERVER_HOST
	serverHost, ok := config["SERVER_HOST"]
	if !ok {
		panic(errors.New("不存在名为SERVER_HOST的配置或配置为空"))
	}
	serverHost_string, ok := serverHost.(string)
	if !ok {
		panic(errors.New("SERVER_HOST必须是字符串类型"))
	}

	// SERVER_PORT
	serverPort, ok := config["SERVER_PORT"]
	if !ok {
		panic(errors.New("不存在名为SERVER_PORT的配置或配置为空"))
	}
	serverPort_int, ok := serverPort.(float64)
	if !ok {
		panic(errors.New("SERVER_PORT必须是int型"))
	}

	// Web_SERVER_PORT
	webServerPort, ok := config["Web_SERVER_PORT"]
	if !ok {
		panic(errors.New("不存在名为Web_SERVER_PORT的配置或配置为空"))
	}
	webServerPort_int, ok := webServerPort.(float64)
	if !ok {
		panic(errors.New("Web_SERVER_PORT必须是int型"))
	}

	// 设置ServerAddress
	ServerAddress = fmt.Sprintf("%s:%d", serverHost_string, int(serverPort_int))

	// 设置WebServerAddress
	WebServerAddress = fmt.Sprintf(":%d", int(webServerPort_int))
}

// 设置日志文件路径
func setLogPath() {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	logPath := filepath.Dir(path)

	logUtil.SetLogPath(filepath.Join(logPath, LOG_PATH_SUFFIX))
}

// 显示数据大小信息
func displayDataSize() {
	for {
		msg := fmt.Sprintf("%s:总共收到%s，发送%s", timeUtil.Format(time.Now(), "yyyy-MM-dd HH:mm:ss"), mathUtil.GetSizeDesc(client.TotalReceiveSize), mathUtil.GetSizeDesc(client.TotalSendSize))
		msg += stringUtil.GetNewLineString()
		msg += fmt.Sprintf("当前客户端数量：%d, 当前玩家数量：%d", len(chatBLL.ClientList), len(chatBLL.PlayerList))
		fmt.Println(msg)
		logUtil.Log(msg, logUtil.Debug, true)
		time.Sleep(time.Minute)
	}
}

// 处理客户端逻辑
// clientObj：客户端对象
// clientAddChan: 客户端增加的channel
// clientRemoveChan: 客户端移除的channel
// playerAddChan: 玩家增加的channel
// playerRemoveChan: 玩家移除的channel
func handleClientContent(clientObj *client.Client, clientAddChan, clientRemoveChan, playerAddChan, playerRemoveChan chan *player.PlayerAndClient) {
	for {
		content, ok := clientObj.GetValieMessage()
		if !ok {
			break
		}

		// 处理数据，如果长度为0则表示心跳包
		if len(content) == 0 {
			continue
		} else {
			chatBLL.HanleRequest(clientObj, content, clientAddChan, clientRemoveChan, playerAddChan, playerRemoveChan)
		}
	}
}

// 处理客户端连接
// conn：客户端连接对象
// clientAddChan: 客户端增加的channel
// clientRemoveChan: 客户端移除的channel
// playerAddChan: 玩家增加的channel
// playerRemoveChan: 玩家移除的channel
func handleConn(conn net.Conn, clientAddChan, clientRemoveChan, playerAddChan, playerRemoveChan chan *player.PlayerAndClient) {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
		}
	}()

	// 创建客户端对象
	clientObj := client.NewClient(&conn, conn)

	// 将客户端对象添加到客户端增加的channel中
	clientAddChan <- player.NewPlayerAndClient(nil, clientObj)

	// 将客户端对象添加到客户端移除的channel中
	defer func() {
		clientRemoveChan <- player.NewPlayerAndClient(nil, clientObj)
	}()

	// 无限循环，不断地读取数据，解析数据，处理数据
	for {
		// 先读取数据，每次读取1024个字节
		readBytes := make([]byte, 1024)
		n, err := conn.Read(readBytes)
		if err != nil {
			var errMsg string

			// 判断是连接关闭错误，还是普通错误
			if err == io.EOF {
				errMsg = fmt.Sprintf("另一端关闭了连接：%s", err)
			} else {
				errMsg = fmt.Sprintf("读取数据错误：%s", err)
			}

			logUtil.Log(errMsg, logUtil.Error, true)

			break
		}

		// 将读取到的数据追加到已获得的数据的末尾
		clientObj.AppendContent(readBytes[:n])

		// 处理数据
		handleClientContent(clientObj, clientAddChan, clientRemoveChan, playerAddChan, playerRemoveChan)
	}
}

// 启动服务器
// ch：用于与main线程传递消息的channel：向ch写入0表示启动服务器成功，1表示失败
func startServer(ch chan int) {
	// 监听指定的端口
	listener, err := net.Listen(SERVER_NETWORK, ServerAddress)
	if err != nil {
		logUtil.Log(fmt.Sprintf("Listen Error: %s", err), logUtil.Error, true)
		ch <- 1
		return
	} else {
		// 写入0表示启动成功，则main线程可以继续往下进行
		logUtil.Log(fmt.Sprintf("Got listener for the server. (local address: %s)", listener.Addr()), logUtil.Debug, true)
		ch <- 0
	}
	defer func() {
		listener.Close()
		ch <- 1
	}()

	for {
		// 阻塞直至新连接到来
		conn, err := listener.Accept()
		if err != nil {
			logUtil.Log(fmt.Sprintf("Accept Error: %s", err), logUtil.Error, true)
			continue
		}

		// 启动一个新协程来处理链接
		go handleConn(conn, chatBLL.ClientAddChan, chatBLL.ClientRemoveChan, chatBLL.PlayerAddChan, chatBLL.PlayerRemoveChan)
	}
}

// 启动web服务器
// ch：用于与main线程传递消息的channel：向ch写入0表示启动服务器成功，2表示失败
func startWebServer(ch chan int) {
	defer func() {
		ch <- 2
	}()

	http.HandleFunc("/", webBLL.ReceiveMessage)
	fmt.Println("WebServerAddress", WebServerAddress)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		logUtil.Log(fmt.Sprintf("ListenAndServer:", err), logUtil.Error, true)
		ch <- 2
	} else {
		// 写入0表示启动成功，则main线程可以继续往下进行
		ch <- 0
		logUtil.Log(fmt.Sprintf("Web server start successfully. (local address: %s)", WebServerAddress), logUtil.Debug, true)
	}
}

func main() {
	ch := make(chan int)

	// 启动服务器
	go startServer(ch)

	// 通过解析从启动服务器的coroutine中返回的值，来判断启动的结果；0表示启动失败，非0表示启动成功
	if <-ch == 1 {
		fmt.Println("Socket服务器启动失败，请检查配置")
		os.Exit(1)
	} else {
		fmt.Println("Socket服务器启动成功，等待客户端的接入。。。")
	}

	// 定时显示数据大小信息
	go displayDataSize()

	// 阻塞，等待输出，以免main线程退出
	<-ch
}
