package znet

import "zinx/ziface"

// 实现 router 时，先嵌入这个基本结构体，然后根据需要对这个基本结构体的方法进行重写
type BaseRouter struct{}

// 这里之所以 BaseRouter 的方法都为空
// 是因为有点 Router 不希望有 PreHandle 或 PostHandle
// 所以 Router 全部嵌入 BaseRouter的好处是，不需要实现 PreHandle 和 PostHandle 也可以被实例化
func (br *BaseRouter) PreHandle(req ziface.IRequest)  {}
func (br *BaseRouter) Handle(req ziface.IRequest)     {}
func (br *BaseRouter) PostHandle(req ziface.IRequest) {}
