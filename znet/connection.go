package znet

import (
	"fmt"
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
	Router ziface.IRouter
	// 告知该连接已经退出/停止的 channel
	ExitBuffChan chan bool
}

// NewConnection 创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool, 1),
	}
}

// 处理 conn 读数据的 goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.Conn.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		// 读取我们最大的数据到 buf 中
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			continue
		}
		// 得到当前客户端请求的 Request 数据
		req := Request{
			conn: c,
			data: buf,
		}
		// 从路由 Routers 中找到注册绑定 Conn的对应 Handle
		go func(request ziface.IRequest) {
			// 执行注册的路由方法
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
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
	if c.isClosed == true {
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
