package ws

import (
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/ports"
	socketio "github.com/googollee/go-socket.io"
)

type userCtx struct {
	ID    int64
	Email string
}

func RegisterSocketHandlers(
	server *socketio.Server,
	hub ports.AdviceHub,
	authService ports.AuthClient,
) {
	handler := onConnectHandler(hub, authService)

	server.OnConnect("/", func(s socketio.Conn) error {
		return handler(&connWrapper{s})
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		if ctx, ok := s.Context().(userCtx); ok {
			hub.Unregister(ctx.ID)
		}
	})
}
