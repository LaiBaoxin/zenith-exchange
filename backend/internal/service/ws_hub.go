package service

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

// SubscriptionRequest 订阅/取消订阅请求格式
type SubscriptionRequest struct {
	Action  string   `json:"action"`  // "subscribe" 或 "unsubscribe"
	Topic   string   `json:"topic"`   // "depth" 或 "trade"
	Symbols []string `json:"symbols"` // ["BTC_USDT"]
}

type Client struct {
	UserID uint64
	Conn   *websocket.Conn
	Send   chan []byte
	// 记录该客户端订阅的房间，用于断开连接时自动清理
	// Key 格式: "depth:BTC_USDT"
	Rooms map[string]bool
	mu    sync.RWMutex
}

type TopicMessage struct {
	Topic   string // "depth" 或 "trade"
	Symbol  string
	Message []byte
}

type Hub struct {
	Clients map[uint64]*Client
	// Rooms 维护 房间名 -> 客户端列表 的映射
	Rooms      map[string]map[*Client]bool
	Broadcast  chan []byte       // 全局广播
	TopicChan  chan TopicMessage // 定点主题广播
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint64]*Client),
		Rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan []byte),
		TopicChan:  make(chan TopicMessage, 100),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				// 清理房间引用
				for roomName := range client.Rooms {
					if room, exists := h.Rooms[roomName]; exists {
						delete(room, client)
					}
				}
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.UserID)
				}
			}
			h.mu.RUnlock()

		case tm := <-h.TopicChan:
			// 定点推送给订阅了该币种的用户
			roomName := tm.Topic + ":" + tm.Symbol
			h.mu.RLock()
			if clients, ok := h.Rooms[roomName]; ok {
				for client := range clients {
					select {
					case client.Send <- tm.Message:
					default:
						// 容错处理
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (c *Client) ReadPump(h *Hub) {
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// 处理订阅逻辑
		var req SubscriptionRequest
		if err := json.Unmarshal(message, &req); err == nil {
			h.handleSubscription(c, req)
		}
	}
}

func (h *Hub) handleSubscription(c *Client, req SubscriptionRequest) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, symbol := range req.Symbols {
		roomName := req.Topic + ":" + symbol
		if req.Action == "subscribe" {
			if h.Rooms[roomName] == nil {
				h.Rooms[roomName] = make(map[*Client]bool)
			}
			h.Rooms[roomName][c] = true
			c.mu.Lock()
			if c.Rooms == nil {
				c.Rooms = make(map[string]bool)
			}
			c.Rooms[roomName] = true
			c.mu.Unlock()
		} else {
			delete(h.Rooms[roomName], c)
		}
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
