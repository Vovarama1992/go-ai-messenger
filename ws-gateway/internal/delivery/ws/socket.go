package ws

import (
	"context"
	"fmt"

	socketio "github.com/googollee/go-socket.io"

	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
)

func RegisterSocketHandlers(
	server *socketio.Server,
	authService ports.AuthService,
	chatService ports.ChatService,
	kafkaProducer ports.KafkaProducer,
	hub *Hub,
) {
	server.OnConnect("/", func(s socketio.Conn) error {
		u := s.URL()
		token := u.Query().Get("token")

		if token == "" {
			return fmt.Errorf("missing token")
		}

		userID, err := authService.ValidateToken(context.Background(), token)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}

		hub.Register(userID, s)
		s.SetContext(userID)
		s.Emit("connected", "✅ WebSocket соединение установлено")
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		if ctx := s.Context(); ctx != nil {
			if userID, ok := ctx.(int64); ok {
				hub.Unregister(userID)
			}
		}
	})

	server.OnEvent("/", "message", func(s socketio.Conn, msg map[string]interface{}) {
		ctx := context.Background()

		userID, ok := s.Context().(int64)
		if !ok {
			s.Emit("error", "unauthorized")
			return
		}

		chatIDFloat, ok := msg["chatId"].(float64)
		if !ok {
			s.Emit("error", "invalid chatId")
			return
		}
		chatID := int64(chatIDFloat)

		text, ok := msg["text"].(string)
		if !ok || text == "" {
			s.Emit("error", "empty message")
			return
		}

		bindings, err := chatService.GetBindingsByChat(ctx, chatID)
		if err != nil {
			s.Emit("error", "internal error")
			return
		}

		payload := map[string]interface{}{
			"chatId":      chatID,
			"senderId":    userID,
			"text":        text,
			"aiGenerated": false,
		}
		_ = kafkaProducer.Produce(ctx, "chat.message.persist", payload)

		for _, b := range bindings {
			if b.Type == "advice" {
				_ = kafkaProducer.Produce(ctx, "chat.message.ai.advice-request", payload)
			}
			if b.Type == "autoreply" {
				replyPayload := map[string]interface{}{
					"chatId":      chatID,
					"senderId":    userID,
					"text":        text,
					"aiGenerated": false,
					"recipientId": b.UserID,
				}
				_ = kafkaProducer.Produce(ctx, "chat.message.ai.autoreply-request", replyPayload)
			}
		}
	})
}
