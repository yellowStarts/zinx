package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

// 连接管理结构体
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接信息
	connLock    sync.RWMutex                  // 读写连接的读写锁
}

// NewConnManager 创建一个连接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加连接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源Map 加 写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 将 conn 连接添加到 ConnManager 中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connection Add to ConnManager successfully: conn num = ", cm.Len())
}

// Remove 删除连接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源Map 加 写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 删除连接信息
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", cm.Len())
}

// Get 利用 connID 获取连接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源Map 加 写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// Len 获取当前连接数量
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// ClearConn 清楚并停止所有连接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源Map 加 写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 停止并删除全部的连接信息
	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, connID)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", cm.Len())
}
