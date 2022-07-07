package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

// IConnection 接口实现，定义一个 Connection 服务结构体
type Connection struct {
	// 当前连接的socker TCP 套接字
	Conn *net.TCPConn
	// 当前连接的ID，也可以称作为 SessionID，ID全局唯一
	ConnID uint32
	// 当前连接的关闭状态
	isClosed bool
	// 该连接的处理方法 router
	MsgHandler ziface.IMsgHandle
	// 告知该连接已经退出/停止的 channel
	ExitBuffChan chan bool
	// 无缓冲管道，用于读，写两个 goroutine 之间的消息通信
	msgChan chan []byte
}

// NewConnection 创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
	}
}

// StartWriter 写消息 Goroutine，用户将数据发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitBuffChan:
			// conn 已经关闭
			return
		}
	}
}

// 处理 conn 读数据的 goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.Conn.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		// 创建拆包解包的对象
		dp := NewDataPack()
		// 读取客户端的 Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitBuffChan <- true
			continue
		}
		// 拆包，得到 msgId 和 dataLen 放在msg 中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- true
			continue
		}
		// 根据 dataLen 读取 data，放在 msg.Data 中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)
		// 得到当前客户端请求的 Request 数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		// 从路由 Routers 中找到注册绑定 Conn的对应 Handle
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

// Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	// 开始处理该连接读取到客户端数据之后的请求业务
	go c.StartReader()
	for {
		select {
		case <-c.ExitBuffChan:
			// 得到退出消息，不在阻塞
			return
		}
	}
}

// Stop 停止连接，结束当前连接状态 M
func (c *Connection) Stop() {
	// 1. 如果当前连接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true
	// TODO Cpnnection Stop() 如果用户注册了该连接的关闭回调业务，那么在此刻应该显示调用
	// 关闭 socket 连接
	c.Conn.Close()
	// 通知从缓冲队列读取数据的业务，该连接已经关闭
	c.ExitBuffChan <- true
	// 关闭该连接全部管道
	close(c.ExitBuffChan)
}

// GetTCPConnection 从当前连接获取原始的 socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 直接将 Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	// 将 data 封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}
	// 写回客户端
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id ", msgId, " error ")
		c.ExitBuffChan <- true
		return errors.New("conn Write error")
	}
	return nil
}
