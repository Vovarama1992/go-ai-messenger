package main

import (
	"log"
	"net/http"
	"os"

	socketio "github.com/googollee/go-socket.io"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/Vovarama1992/go-ai-messenger/proto/authpb"
	chatpb "github.com/Vovarama1992/go-ai-messenger/proto/chatpb"

	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/delivery/ws"
	grpcadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/grpc"
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
)

func main() {
	port := os.Getenv("WS_GATEWAY_PORT")
	if port == "" {
		log.Fatal("WS_GATEWAY_PORT is not set")
	}

	authAddr := os.Getenv("AUTH_SERVICE_GRPC_ADDR")
	if authAddr == "" {
		log.Fatal("AUTH_SERVICE_GRPC_ADDR is not set")
	}

	authConn, err := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to auth-service: %v", err)
	}
	defer authConn.Close()
	authGrpcClient := authpb.NewAuthServiceClient(authConn)
	authService := grpcadapter.NewAuthService(authGrpcClient)

	chatAddr := os.Getenv("CHAT_SERVICE_GRPC_ADDR")
	if chatAddr == "" {
		log.Fatal("CHAT_SERVICE_GRPC_ADDR is not set")
	}

	chatConn, err := grpc.Dial(chatAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to chat-service: %v", err)
	}
	defer chatConn.Close()
	chatClient := chatpb.NewChatServiceClient(chatConn)
	chatService := grpcadapter.NewChatService(chatClient)

	producer := kafkaadapter.NewKafkaProducer()
	if err := producer.Start(); err != nil {
		log.Fatalf("❌ Failed to start Kafka producer: %v", err)
	}
	defer func() {
		if err := producer.Stop(); err != nil {
			log.Printf("❌ Failed to stop Kafka producer: %v", err)
		}
	}()

	hub := ws.NewHub()

	server := socketio.NewServer(nil)
	ws.RegisterSocketHandlers(server, authService, chatService, producer, hub)

	listener := ws.NewForwardListener(hub)
	listener.Start()
	defer listener.Stop()

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)

	log.Printf("🚀 WebSocket-сервер запущен на :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
