package ws

import (
	"context"
	"fmt"
	"os"

	socketio "github.com/googollee/go-socket.io"

	adapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/adapters/ws"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
)

type UserCtx struct {
	ID    int64
	Email string
}

func RegisterSocketHandlers(
	server *socketio.Server,
	authService ports.AuthService,
	chatService ports.ChatService,
	kafkaProducer ports.KafkaProducer,
	hub ports.Hub,
) {
	topicPersist := os.Getenv("TOPIC_MESSAGE_PERSIST")
	topicFeed := os.Getenv("TOPIC_AI_FEED")

	if topicPersist == "" || topicFeed == "" {
		panic("❌ Required Kafka topic envs are not set")
	}

	server.OnConnect("/", MakeConnectHandler(authService, hub))
	server.OnDisconnect("/", MakeDisconnectHandler(hub))
	server.OnEvent("/", "message", MakeMessageHandler(chatService, kafkaProducer, topicPersist, topicFeed))
	server.OnEvent("/", "join-room", MakeJoinRoomHandler(hub))
}

func MakeConnectHandler(auth ports.AuthService, hub ports.Hub) func(socketio.Conn) error {
	return func(s socketio.Conn) error {
		u := s.URL()
		token := u.Query().Get("token")
		if token == "" {
			return fmt.Errorf("missing token")
		}
		conn := adapter.NewSocketConnAdapter(s)
		return HandleConnect(auth, hub, conn, token)
	}
}

func MakeDisconnectHandler(hub ports.Hub) func(socketio.Conn, string) {
	return func(s socketio.Conn, _ string) {
		if ctx, ok := s.Context().(UserCtx); ok {
			hub.Unregister(ctx.ID)
		}
	}
}

func MakeMessageHandler(
	chatService ports.ChatService,
	kafka ports.KafkaProducer,
	topicPersist string,
	topicFeed string,
) func(socketio.Conn, map[string]interface{}) {
	return func(s socketio.Conn, msg map[string]interface{}) {
		ctx := context.Background()

		user, ok := s.Context().(UserCtx)
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

		kafka.Produce(ctx, topicPersist, map[string]interface{}{
			"chatId":      chatID,
			"senderId":    user.ID,
			"text":        text,
			"aiGenerated": false,
		})

		for _, b := range bindings {
			kafka.Produce(ctx, topicFeed, map[string]interface{}{
				"senderEmail": user.Email,
				"text":        text,
				"threadId":    b.ThreadID,
				"bindingType": b.BindingType,
			})
		}
	}
}

func MakeJoinRoomHandler(hub ports.Hub) func(socketio.Conn, map[string]interface{}) {
	return func(s socketio.Conn, msg map[string]interface{}) {
		user, ok := s.Context().(UserCtx)
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
	}
}

func HandleConnect(auth ports.AuthService, hub ports.Hub, conn ports.Conn, rawToken string) error {
	userID, email, err := auth.ValidateToken(context.Background(), rawToken)
	if err != nil {
		return fmt.Errorf("unauthorized")
	}

	hub.Register(userID, conn)
	conn.SetContext(UserCtx{ID: userID, Email: email})
	conn.Emit("connected", "✅ WebSocket соединение установлено")
	return nil
}
