package ws

import (
	"sync"

	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
)

type Hub struct {
	mu      sync.RWMutex
	sockets map[int64]ports.Conn         // userID → socket.Conn
	rooms   map[int64]map[int64]struct{} // chatID → set of userIDs
}

func NewHub() *Hub {
	return &Hub{
		sockets: make(map[int64]ports.Conn),
		rooms:   make(map[int64]map[int64]struct{}),
	}
}

func (h *Hub) Register(userID int64, conn ports.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sockets[userID] = conn
}

func (h *Hub) Unregister(userID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.sockets, userID)

	for chatID, users := range h.rooms {
		delete(users, userID)
		if len(users) == 0 {
			delete(h.rooms, chatID)
		}
	}
}

func (h *Hub) JoinRoom(userID int64, chatID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[chatID]; !ok {
		h.rooms[chatID] = make(map[int64]struct{})
	}
	h.rooms[chatID][userID] = struct{}{}
}

func (h *Hub) LeaveRoom(userID int64, chatID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if users, ok := h.rooms[chatID]; ok {
		delete(users, userID)
		if len(users) == 0 {
			delete(h.rooms, chatID)
		}
	}
}

func (h *Hub) SendToRoom(chatID int64, event string, data any) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users, ok := h.rooms[chatID]
	if !ok {
		return
	}

	for userID := range users {
		if conn, ok := h.sockets[userID]; ok {
			conn.Emit(event, data)
		}
	}
}

func (h *Hub) HasConnection(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.sockets[userID]
	return exists
}

func (h *Hub) GetConn(userID int64) ports.Conn {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conn, _ := h.sockets[userID]
	return conn
}
