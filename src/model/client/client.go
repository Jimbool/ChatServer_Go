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
	ByterOrder = binary.LittleEndian

	// 收到的总数据大小，以B为单位
	TotalReceiveSize int64

	// 发送的总数据大小，以B为单位
	TotalSendSize int64
)

// 定义客户端对象，以实现对客户端连接的封装
type Client struct {
	// 公共属性
	// 唯一标识
	Id *net.Conn

	// 客户端连接对象
	Conn net.Conn

	// 玩家Id
	PlayerId string

	// 私有属性，内部使用
	// 接收到的消息内容
	content []byte

	// 上次活跃时间
	activeTime time.Time
}

// 新建客户端对象
// id：连接对象的指针
// conn：连接对象
// 返回值：客户端对象的指针
func NewClient(id *net.Conn, conn net.Conn) *Client {
	return &Client{
		// 公共属性赋值
		Id:       id,
		Conn:     conn,
		PlayerId: "",

		// 私有属性赋值
		content:    make([]byte, 0, 1024),
		activeTime: time.Now(),
	}
}

// 追加内容
// content：新的内容
// 返回值：无
func (c *Client) AppendContent(content []byte) {
	c.content = append(c.content, content...)
	c.activeTime = time.Now()

	// 增加接收数据量
	atomic.AddInt64(&TotalReceiveSize, int64(len(content)))
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
	contentLength := intAndBytesUtil.BytesToInt(header, ByterOrder)

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
	header := intAndBytesUtil.IntToBytes(contentLength, ByterOrder)

	// 将头部与内容组合在一起
	message := append(header, b...)

	// 增加发送量(包括包头的长度+内容的长度)
	atomic.AddInt64(&TotalSendSize, int64(HEADER_LENGTH+contentLength))

	// 发送消息
	c.Conn.Write(message)
}

// 判断客户端是否超时
// 返回值：是否超时
func (c *Client) HasExpired() bool {
	return c.activeTime.Add(configBLL.ClientExpiredSeconds*time.Second).Unix() < time.Now().Unix()
}

// 玩家登陆
// playerId：玩家Id
// 返回值：无
func (c *Client) PlayerLogin(playerId string) {
	c.PlayerId = playerId
}

// 玩家登出
// 返回值：无
func (c *Client) PlayerLogout() {
	c.PlayerId = ""
}

// 退出
// 返回值：无
func (c *Client) Quit() {
	c.Conn.Close()
}
