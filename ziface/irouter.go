package ziface

// 路由接口
// 这里面路由是 使用框架者给该连接自定的 处理业务方法
// 路由里的 IRequest 则包含用该连接的连接信息和该连接的请求数据信息

type IRouter interface {
	PreHandle(request IRequest)  // 在处理 conn 业务之前的钩子方法
	Handle(request IRequest)     // 处理 conn 业务的方法
	PostHandle(request IRequest) // 处理 conn 业务之后的钩子方法
}
