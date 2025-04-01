package utils

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WebSocketHook struct {
	clients      map[*websocket.Conn]bool // 存储 WebSocket 客户端连接
	mu           sync.Mutex               // 保护 clients 的并发访问
	broadcast    chan string              // 用于广播日志消息的通道
	logHistory   []string                 // 保存日志历史
	historyLimit int                      // 限制日志历史的最大条数
}

// 创建 WebSocketHook 实例
func NewWebSocketHook() *WebSocketHook {
	hook := &WebSocketHook{
		clients:      make(map[*websocket.Conn]bool),
		broadcast:    make(chan string),
		logHistory:   []string{},
		historyLimit: 300, // 设置日志历史的最大条数
	}
	go hook.handleBroadcast()
	return hook
}

// Fire 方法：当 logrus 记录日志时触发
func (hook *WebSocketHook) Fire(entry *logrus.Entry) error {
	// 将日志格式化为字符串
	msg, err := entry.String()
	if err != nil {
		return err
	}

	// 保存到日志历史
	hook.mu.Lock()

	// 使用 termenv 将 ANSI 转义序列转换为 HTML
	hook.logHistory = append(hook.logHistory, msg)
	if len(hook.logHistory) > hook.historyLimit {
		hook.logHistory = hook.logHistory[1:] // 删除最早的日志，保持限制
	}
	hook.mu.Unlock()

	// 将日志消息发送到广播通道
	hook.broadcast <- msg
	return nil
}

// Levels 方法：定义 Hook 监听的日志级别
func (hook *WebSocketHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// 添加 WebSocket 客户端
func (hook *WebSocketHook) AddClient(conn *websocket.Conn) {
	hook.mu.Lock()
	defer hook.mu.Unlock()

	// 将客户端加入到 clients
	hook.clients[conn] = true

	history := map[string][]string{"history":{}}

	history["history"] = hook.logHistory

	jsonHistory, err := json.Marshal(history)
	if err != nil {
		return
	}

	// 发送历史日志给新连接的客户端
	conn.WriteMessage(websocket.TextMessage, []byte(jsonHistory))
}

// 移除 WebSocket 客户端
func (hook *WebSocketHook) RemoveClient(conn *websocket.Conn) {
	hook.mu.Lock()
	defer hook.mu.Unlock()
	delete(hook.clients, conn)
	conn.Close()
}

// 处理广播通道：将日志消息发送到所有客户端
func (hook *WebSocketHook) handleBroadcast() {
	for m := range hook.broadcast {
		hook.mu.Lock()
		msg := map[string]string{"msg": m}
		jsonMsg, err := json.Marshal(msg)    // 使用 json.Marshal 转换为合法的 JSON 字符串
		if err != nil {
			hook.mu.Unlock()
			return
		}
		for client := range hook.clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(jsonMsg))
			if err != nil {
				// 如果发送失败，移除客户端
				hook.RemoveClient(client)
			}
		}
		hook.mu.Unlock()
	}
}
