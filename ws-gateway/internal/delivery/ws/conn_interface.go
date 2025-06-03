package ws

import socketio "github.com/googollee/go-socket.io"

type Conn interface {
	socketio.Conn
}
