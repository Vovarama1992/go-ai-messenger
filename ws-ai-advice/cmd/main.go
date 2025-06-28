package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	socketio "github.com/googollee/go-socket.io"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/Vovarama1992/go-ai-messenger/proto/authpb"
	chatpb "github.com/Vovarama1992/go-ai-messenger/proto/chatpb"

	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/app"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/delivery/ws"
	grpcadapter "github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/infra/grpc"
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/infra/kafka"
)

func main() {
	port := os.Getenv("WS_AI_ADVICE_PORT")
	if port == "" {
		log.Fatal("WS_AI_ADVICE_PORT is not set")
	}

	topic := os.Getenv("TOPIC_AI_ADVICE_RESPONSE")
	if topic == "" {
		log.Fatal("TOPIC_AI_ADVICE_RESPONSE is not set")
	}

	// gRPC connections
	authConn, err := grpc.Dial(os.Getenv("AUTH_GRPC_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ failed to connect to auth gRPC: %v", err)
	}
	defer authConn.Close()

	chatConn, err := grpc.Dial(os.Getenv("CHAT_GRPC_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ failed to connect to chat gRPC: %v", err)
	}
	defer chatConn.Close()

	// Services
	authService := grpcadapter.NewAuthService(authpb.NewAuthServiceClient(authConn))
	chatService := grpcadapter.NewChatService(chatpb.NewChatServiceClient(chatConn))

	// WebSocket setup
	hub := ws.NewHub()
	server := socketio.NewServer(nil)
	ws.RegisterSocketHandlers(server, hub, authService)

	// HTTP WebSocket server
	go func() {
		http.Handle("/socket.io/", server)
		log.Printf("🚀 ws-ai-advice up on :%s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("❌ ListenAndServe: %v", err)
		}
	}()

	// Graceful shutdown context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Pipeline
	go app.RunChannelsBetweener(ctx, chatService)
	go app.RunAdvicePusherToFronts(ctx, hub)

	// Kafka consumer
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092"
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokers},
		GroupID: "ws-ai-advice",
		Topic:   topic,
	})
	consumer := kafkaadapter.NewAdviceConsumer(reader)

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Fatalf("❌ consumer start failed: %v", err)
		}
	}()

	// OS signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("🛑 shutdown initiated")
	consumer.Close()
	server.Close()
}
