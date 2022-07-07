package znet

import "zinx/ziface"

type Request struct {
	// 已经和客户端建立好的 连接
	conn ziface.IConnection
	// 客户端请求的数据
	data []byte
}

// 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.data
}
