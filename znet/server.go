package znet

import (
	"errors"
	"fmt"
	"net"
	"time"
	"zinx/utils"
	"zinx/ziface"
)

// IServer 接口实现，定义一个 Server 服务结构体
type Server struct {
	// 服务器名称
	Name string
	// tcp4 or other
	IPVersion string
	// 服务绑定的 IP 地址
	IP string
	// 服务绑定的端口
	Port int
	// 当前 Server 由用户绑定的回调router，也就是Server注册的连接对应的处理业务
	msgHandler ziface.IMsgHandle
	// 当前Server的连接管理器
	ConnMgr ziface.IConnManager
	// 两个hook函数原型
	// 该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	// 该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

// 定义当前客户端连接的 handle api
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	// 回显业务
	fmt.Println("[Conn Handle] CallBackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// NewServer 创建一个服务器句柄
func NewServer() ziface.IServer {
	// 先初始化全局配置文件
	utils.GlobalObject.Reload()
	s := &Server{
		Name:       utils.GlobalObject.Name, // 从全局参数获取
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,    // 从全局参数获取
		Port:       utils.GlobalObject.TcpPort, // 从全局参数获取
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(), // 创建 ConnManager
	}
	return s
}

// Start 开启网络服务
func (s *Server) Start() {
	// 如果方法正确，运行时就会打印出
	// Version: V0.4, MaxConn: 3, MaxPacketSize: 4096
	fmt.Printf("[START] Server name: %s, listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	// 开启一个 go 去做服务端 Listenner 业务
	go func() {
		// 0. 启动 worker 工作池机制
		s.msgHandler.StartWorkerPool()
		// 1. 获取一个 TCP 的 Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		// 2. 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err", err)
			return
		}
		// 已经监听成功
		fmt.Println("start Zinx server ", s.Name, " success, now listenning...")
		// TODO server.go 应该有一个自动生成 ID 的方法
		var cid uint32
		cid = 0
		// 3. 启动 server 网络链接业务
		for {
			// 3.1 阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			// 3.2 TODO Server.Start() 设置服务器最大连接控制，超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}
			// 3.3 TODO Server.Start() 处理该新连接请求的 业务 方法，此时应该有 handler 和 conn 是绑定的
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++
			// 3.4 启动当前连接的处理业务
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] zinx server, name ", s.Name)
	// 将其他需要清理的连接信息或者其他信息，也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

// Serve 开启服务
func (s *Server) Serve() {
	s.Start()
	// TODO Server.Serve() 是否在启动服务的时候，还要处理其他的事情呢，可以在这里添加
	// 阻塞，否则主 go 退出，listenner 的 go 将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}

// AddRouter 路由功能
// 给当前服务注册一个路由业务方法，供客户端连接处理使用
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

// GetConnMgr 得到连接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接 OnConnStart hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接 OnConnStop hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
