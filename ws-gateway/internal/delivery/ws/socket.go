package ws

import (
	"context"
	"fmt"
	"os"

	socketio "github.com/googollee/go-socket.io"

	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
)

type userCtx struct {
	ID    int64
	Email string
}

func RegisterSocketHandlers(
	server *socketio.Server,
	authService ports.AuthService,
	chatService ports.ChatService,
	kafkaProducer ports.KafkaProducer,
	hub *Hub,
) {
	// читаем топики из ENV
	topicPersist := os.Getenv("TOPIC_MESSAGE_PERSIST")
	topicFeed := os.Getenv("TOPIC_AI_FEED")

	if topicPersist == "" || topicFeed == "" {
		panic("❌ Required Kafka topic envs are not set")
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		u := s.URL()
		token := u.Query().Get("token")
		if token == "" {
			return fmt.Errorf("missing token")
		}

		userID, email, err := authService.ValidateToken(context.Background(), token)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}

		hub.Register(userID, s)
		s.SetContext(userCtx{ID: userID, Email: email})
		s.Emit("connected", "✅ WebSocket соединение установлено")
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		if ctx, ok := s.Context().(userCtx); ok {
			hub.Unregister(ctx.ID)
		}
	})

	server.OnEvent("/", "message", func(s socketio.Conn, msg map[string]interface{}) {
		ctx := context.Background()

		user, ok := s.Context().(userCtx)
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

		// persist → ID
		persistPayload := map[string]interface{}{
			"chatId":      chatID,
			"senderId":    user.ID,
			"text":        text,
			"aiGenerated": false,
		}
		_ = kafkaProducer.Produce(ctx, topicPersist, persistPayload)

		// feed в AI
		for _, b := range bindings {
			_ = kafkaProducer.Produce(ctx, topicFeed, map[string]interface{}{
				"senderEmail": user.Email,
				"text":        text,
				"threadId":    b.ThreadID,
				"bindingType": b.Type,
			})
		}
	})

	server.OnEvent("/", "join-room", func(s socketio.Conn, msg map[string]interface{}) {

		user, ok := s.Context().(userCtx)
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

		hub.JoinRoom(user.ID, chatID)
		s.Emit("joined-room", fmt.Sprintf("Subscribed to chat %d", chatID))
	})
}
