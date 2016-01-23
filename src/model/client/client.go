/*
客户端对象包，将服务器与客户端的连接net.Conn封装起来，并进行管理
*/
package client

import (
	"encoding/binary"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"github.com/Jordanzuo/goutil/intAndBytesUtil"
	"net"
	"sync/atomic"
	"time"
)

const (
	// 包头的长度
	HEADER_LENGTH = 4
)

var (
	// 字节的大小端顺序
	byterOrder = binary.LittleEndian

	// 收到的总数据大小，以B为单位
	totalReceiveSize int64

	// 发送的总数据大小，以B为单位
	totalSendSize int64

	// 全局客户端的id，从1开始进行自增
	globalClientId int32 = 0
)

// 获得自增的id值
func getIncrementId() int32 {
	atomic.AddInt32(&globalClientId, 1)

	return globalClientId
}

// 获取收到的总数据大小，以B为单位
// 返回值：
// 收到的总数据大小，以B为单位
func TotalReceiveSize() int64 {
	return totalReceiveSize
}

// 获取发送的总数据大小，以B为单位
// 返回值：
// 发送的总数据大小，以B为单位
func TotalSendSize() int64 {
	return totalSendSize
}

// 定义客户端对象，以实现对客户端连接的封装
type Client struct {
	// 唯一标识
	id int32

	// 客户端连接对象
	conn net.Conn

	// 玩家Id
	playerId string

	// 接收到的消息内容
	content []byte

	// 上次活跃时间
	activeTime time.Time
}

// 新建客户端对象
// conn：连接对象
// 返回值：客户端对象的指针
func NewClient(conn net.Conn) *Client {
	return &Client{
		id:         getIncrementId(),
		conn:       conn,
		playerId:   "",
		content:    make([]byte, 0, 1024),
		activeTime: time.Now(),
	}
}

// 获取唯一标识
func (c *Client) Id() int32 {
	return c.id
}

// 获取玩家Id
// 返回值：
// 玩家Id
func (c *Client) PlayerId() string {
	return c.playerId
}

// 追加内容
// content：新的内容
// 返回值：无
func (c *Client) AppendContent(content []byte) {
	c.content = append(c.content, content...)
	c.activeTime = time.Now()

	// 增加接收数据量
	atomic.AddInt64(&totalReceiveSize, int64(len(content)))
}

// 获取有效的消息
// 返回值：消息内容
//		：是否含有有效数据
func (c *Client) GetValieMessage() ([]byte, bool) {
	// 判断是否包含头部信息
	if len(c.content) < HEADER_LENGTH {
		return nil, false
	}

	// 获取头部信息
	header := c.content[:HEADER_LENGTH]

	// 将头部数据转换为内部的长度
	contentLength := intAndBytesUtil.BytesToInt(header, byterOrder)

	// 判断长度是否满足
	if len(c.content)-HEADER_LENGTH < contentLength {
		return nil, false
	}

	// 提取消息内容
	content := c.content[HEADER_LENGTH : HEADER_LENGTH+contentLength]

	// 将对应的数据截断，以得到新的数据
	c.content = c.content[HEADER_LENGTH+contentLength:]

	return content, true
}

// 发送字节数组消息
// b：待发送的字节数组
func (c *Client) SendByteMessage(b []byte) {
	// 获得数组的长度
	contentLength := len(b)

	// 将长度转化为字节数组
	header := intAndBytesUtil.IntToBytes(contentLength, byterOrder)

	// 将头部与内容组合在一起
	message := append(header, b...)

	// 增加发送量(包括包头的长度+内容的长度)
	atomic.AddInt64(&totalSendSize, int64(HEADER_LENGTH+contentLength))

	// 发送消息
	c.conn.Write(message)
}

// 判断客户端是否超时
// 返回值：是否超时
func (c *Client) HasExpired() bool {
	return time.Now().Unix() > c.activeTime.Add(configBLL.ClientExpireSeconds()*time.Second).Unix()
}

// 玩家登陆
// playerId：玩家Id
// 返回值：无
func (c *Client) PlayerLogin(playerId string) {
	c.playerId = playerId
}

// 玩家登出
// 返回值：无
func (c *Client) PlayerLogout() {
	c.playerId = ""
}

// 玩家登出，客户端退出
// 返回值：无
func (c *Client) LogoutAndQuit() {
	c.PlayerLogout()
	c.conn.Close()
}
