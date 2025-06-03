package ws

import (
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

type Hub struct {
	mu      sync.RWMutex
	sockets map[int64]socketio.Conn // userID → socket.Conn
}

func NewHub() *Hub {
	return &Hub{
		sockets: make(map[int64]socketio.Conn),
	}
}

func (h *Hub) Register(userID int64, conn socketio.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sockets[userID] = conn
}

func (h *Hub) Unregister(userID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.sockets, userID)
}

func (h *Hub) Send(userID int64, event string, data any) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if conn, ok := h.sockets[userID]; ok {
		conn.Emit(event, data)
	}
}
