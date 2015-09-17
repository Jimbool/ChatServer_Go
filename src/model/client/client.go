package client

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Jordanzuo/goutil/intAndBytesUtil"
	"github.com/Jordanzuo/goutil/logUtil"
	"net"
	"time"
)

const (
	// 包头的长度
	HEADER_LENGTH = 4
)

var (
	// 字节的大小端顺序
	ByterOrder = binary.LittleEndian

	// 客户端过期的秒数
	ClientExpiredSeconds time.Duration
)

// 设置参数
// config：从配置文件里面解析出来的配置内容
func SetParam(config map[string]interface{}) {
	// 解析ClientExpiredSeconds
	clientExpiredSeconds, ok := config["CLIENT_EXPIRED_SECONDS"]
	if !ok {
		panic(errors.New("不存在名为CLIENT_EXPIRED_SECONDS的配置或配置为空"))
	}
	clientExpiredSeconds_int, ok := clientExpiredSeconds.(float64)
	if !ok {
		panic(errors.New("CLIENT_EXPIRED_SECONDS必须是int型"))
	}

	// 设置client的参数：ClientExpiredSeconds
	ClientExpiredSeconds = time.Duration(int(clientExpiredSeconds_int))
}

// 定义客户端对象，以实现对客户端连接的封装
type Client struct {
	// 公共属性
	// 唯一标识
	Id *net.Conn

	// 客户端连接对象
	Conn net.Conn

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
		Id:   id,
		Conn: conn,

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
}

// 判断当前已经收到的消息是否为有效消息（即已经收到了完整的信息）
// 返回值：是否有效
func (c *Client) ifHasValidMessage() bool {
	// 判断是否包含头部信息
	if len(c.content) < HEADER_LENGTH {
		return false
	}

	// 获取头部信息
	header := c.content[:HEADER_LENGTH]

	// 将头部数据转换为内部的长度
	contentLength := intAndBytesUtil.BytesToInt(header, ByterOrder)

	// 判断长度是否满足
	if len(c.content)-HEADER_LENGTH < contentLength {
		return false
	}

	return true
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

	// 发送消息
	_, err := c.Conn.Write(message)
	if err != nil {
		logUtil.Log(fmt.Sprintf("发送数据错误：%s", err), logUtil.Error, true)
	}
}

// 判断客户端是否超时
// 返回值：是否超时
func (c *Client) IfExpired() bool {
	return c.activeTime.Add(ClientExpiredSeconds*time.Second).Unix() < time.Now().Unix()
}
