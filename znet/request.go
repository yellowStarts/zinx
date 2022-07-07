package znet

import "zinx/ziface"

// 请求结构体
type Request struct {
	// 已经和客户端建立好的 连接
	conn ziface.IConnection
	// 客户端请求的数据
	msg ziface.IMessage
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgId 获取请求的消息的ID
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
