package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter // 一定要先嵌入基础 BaseRouter
}

// Test Handle
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// Server 模块的测试函数
func main() {
	// 1 创建一个 server 句柄 s
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	// 2. 开启服务
	s.Serve()
}
