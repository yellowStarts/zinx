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

// Test PreHandle
func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Before ping ...\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// Test Handle
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// Test PostHandle
func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping ...\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// Server 模块的测试函数
func main() {
	// 1 创建一个 server 句柄 s
	s := znet.NewServer("[zinx V0.3]")
	s.AddRouter(&PingRouter{})
	// 2. 开启服务
	s.Serve()
}
