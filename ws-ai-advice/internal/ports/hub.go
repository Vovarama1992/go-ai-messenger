package ports

import socketio "github.com/googollee/go-socket.io"

type AdviceHub interface {
	Register(userID int64, conn socketio.Conn)
	Unregister(userID int64)
	Send(userID int64, event string, data any)
}
