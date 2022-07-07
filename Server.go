package main

import (
	"fmt"
	"io"
	"net"
	"zinx/znet"
)

// Server 模块的测试函数
// 知识负责测试 datapack 拆包，封包功能
func main() {
	// 常见socket TCP server
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}
	// 创建服务器goroutine，负责从客户端goroutine读取粘包的数据，然后进行解析
	for {
		conn, err := listenner.Accept()
		if err != nil {
			fmt.Println("server accept err: ", err)
		}
		// 处理客户端请求
		go func(conn net.Conn) {
			// 创建封包拆包对象
			dp := znet.NewDataPack()
			for {
				// 1. 先读出流中的 head 部分
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData) // ReadFull 会把 msg 填充满为止
				if err != nil {
					fmt.Println("read head error")
					break
				}
				// 将 headData 字节流 拆包到 msg 中
				msgHead, err := dp.UnPack(headData)
				if err != nil {
					fmt.Println("server unpack err: ", err)
					return
				}
				if msgHead.GetDataLen() > 0 {
					// msg 是有 data 数据，需要再次读取 data 数据
					msg := msgHead.(*znet.Message)
					msg.Data = make([]byte, msg.GetDataLen())
					// 根据 dataLen 从io中读取字节流
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err: ", err)
						return
					}
					fmt.Println("==> Recv Msg: ID=", msg.Id, " , len=", msg.DataLen, " ,data=", string(msg.Data))
				}
			}
		}(conn)
	}
}
