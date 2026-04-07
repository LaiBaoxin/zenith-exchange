package service

import (
	"github.com/gorilla/websocket"
)

// Client 代表一个活跃的 WS 连接
type Client struct {
	UserID uint64
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub 管理所有客户端并分发消息
type Hub struct {
	// 注册：key 为 UserID
	Clients    map[uint64]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint64]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
		// 当有全局广播（如最新成交价）时
		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				client.Send <- message
			}
		}
	}
}
