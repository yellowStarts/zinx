package ziface

import "net"

// 定义连接接口
type IConnection interface {
	// 启动连接，让当前连接开始工作
	Start()
	// 停止连接，结束当前连接状态M
	Stop()
	// 从当前连接获取原始的 socket TCPConn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接 ID
	GetConnID() uint32
	// 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// 直接将 Message 数据发送给远程的TCP客户端
	SendMsg(msgId uint32, data []byte) error
	// 直接将Message数据发送给远程的TCP客户端(有缓冲)
	SendBuffMsg(msgId uint32, data []byte) error
	// 设置链接属性
	SetProperty(key string, value interface{})
	// 获取链接属性
	GetProperty(key string) (interface{}, error)
	// 移除链接属性
	RemoveProperty(key string)
}

// 定义一个统一处理连接业务员的接口
// 第一参数是 socket 原生链接
// 第二个参数是客户端请求的数据
// 第三个参数是客户端请求的数据长度
type HandFunc func(*net.TCPConn, []byte, int) error
