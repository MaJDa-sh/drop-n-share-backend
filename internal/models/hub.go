package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	rooms map[string]*Room
	mu    sync.Mutex
}

type Room struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

func (h *Hub) GetRoom(roomID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, exists := h.rooms[roomID]; exists {
		return room
	}

	room := &Room{clients: make(map[*websocket.Conn]bool)}
	h.rooms[roomID] = room
	return room
}

func (r *Room) RegisterClient(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[conn] = true
}

func (r *Room) RemoveClient(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, conn)
}

func (r *Room) BroadcastMessage(message []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for client := range r.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			client.Close()
			r.RemoveClient(client)
		}
	}
}

func (r *Room) GetClientList() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	clientList := make([]string, 0, len(r.clients))
	for client := range r.clients {
		clientList = append(clientList, client.RemoteAddr().String())
	}
	return clientList
}
