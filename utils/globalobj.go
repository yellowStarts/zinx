package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

// 存储一切有关 zinx 框架的全局参数，供其他模块使用
// 一些参数也可以通过 用户根据 zinx.json 来配置
type GlobalObj struct {
	TcpServer        ziface.IServer // 当前 zinx 的全服 server 对象
	Host             string         // 当前服务器主机 IP
	TcpPort          int            // 当前服务器主机监听端口
	Name             string         // 当前服务器名称
	Version          string         // 当前 zinx 版本
	MaxPacketSize    uint32         // 都需数据包的最大值
	MaxConn          int            // 当前服务器主机允许的最大连接个数
	WorkerPoolSize   uint32         // 业务工作Worker池的数量
	MaxWorkerTaskLen uint32         // 业务工作Worker对应负责的任务队列最大任务存储数量
	ConfFilePath     string         // 配置文件路径
}

// 定义一个全局对象
var GlobalObject *GlobalObj

// Reload 读取用户的配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	// 将 json 数据解析到 struct 中
	// fmt.Printf("json :%s", data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供 init 方法，默认加载
func init() {
	// 初始化 GlobalObject 变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		TcpPort:          7777,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     "conf/zinx.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	// 从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
