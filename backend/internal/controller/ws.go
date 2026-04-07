package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // 开发环境允许跨域
}

type WSHandler struct {
	hub *service.Hub
}

func NewWSHandler(h *service.Hub) *WSHandler {
	return &WSHandler{hub: h}
}

func (h *WSHandler) HandleWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	userID := c.GetUint64("user_id")

	client := &service.Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.hub.Register <- client

	// 开启两个协程：读和写
	go client.WritePump()
	go client.ReadPump(h.hub)
}
