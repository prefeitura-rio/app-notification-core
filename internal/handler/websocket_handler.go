package handler

import (
	"net/http"

	"github.com/prefeitura-rio/app-notification-core/internal/websocket"
	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub *websocket.Hub
}

func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// ServeWS godoc
// @Summary Conectar ao WebSocket
// @Description Estabelece conexão WebSocket para receber notificações em tempo real
// @Tags websocket
// @Param user_id query string true "Identificador do usuário (CPF, telefone ou email)"
// @Success 101 {string} string "Switching Protocols - WebSocket established"
// @Failure 400 {object} map[string]string
// @Router /ws [get]
func (h *WebSocketHandler) ServeWS(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := websocket.NewClient(h.hub, conn, userID)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
