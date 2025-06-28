package ws

import (
	socketio "github.com/googollee/go-socket.io"
)

type connWrapper struct {
	socketio.Conn
}

func (c *connWrapper) GetToken() string {
	if c == nil {
		return ""
	}
	u := c.URL()
	return u.Query().Get("token")
}
