package websocket

import (
	"encoding/json"
	"github.com/fzolio/app-notification-core/internal/entity"
	"sync"
)

type Hub struct {
	clients    map[string]map[*Client]bool
	broadcast  chan *entity.Notification
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan *entity.Notification, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*Client]bool)
			}
			h.clients[client.userID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.userID)
					}
				}
			}
			h.mu.Unlock()

		case notification := <-h.broadcast:
			h.sendNotification(notification)
		}
	}
}

func (h *Hub) sendNotification(notification *entity.Notification) {
	message, err := json.Marshal(notification)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if notification.Broadcast {
		for _, clients := range h.clients {
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients[client.userID], client)
				}
			}
		}
		return
	}

	var targetUserID string
	if notification.UserCPF != nil {
		targetUserID = *notification.UserCPF
	} else if notification.UserPhone != nil {
		targetUserID = *notification.UserPhone
	}

	if targetUserID != "" {
		if clients, ok := h.clients[targetUserID]; ok {
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastNotification(notification *entity.Notification) {
	h.broadcast <- notification
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}
