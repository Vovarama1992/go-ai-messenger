package ws

import (
	"context"
	"errors"
	"fmt"

	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/ports"
)

func onConnectHandler(hub ports.AdviceHub, auth ports.AuthClient) func(ports.Conn) error {
	return func(conn ports.Conn) error {
		token := conn.GetToken()
		if token == "" {
			return errors.New("missing token")
		}

		userID, email, err := auth.ValidateToken(context.Background(), token)
		if err != nil {
			return fmt.Errorf("unauthorized: %w", err)
		}

		hub.Register(userID, conn)
		conn.SetContext(userCtx{ID: userID, Email: email})
		conn.Emit("connected", "✅ WebSocket соединение установлено")
		return nil
	}
}

func TestableOnConnectHandler(
	hub ports.AdviceHub,
	authService ports.AuthClient,
) func(ports.Conn) error {
	return onConnectHandler(hub, authService)
}
