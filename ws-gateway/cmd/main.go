package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

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

	forwardTopic := os.Getenv("TOPIC_FORWARD_MESSAGE")
	if forwardTopic == "" {
		log.Fatal("TOPIC_FORWARD_MESSAGE is not set")
	}

	inviteTopic := os.Getenv("TOPIC_CHAT_INVITE")
	if inviteTopic == "" {
		log.Fatal("TOPIC_CHAT_INVITE is not set")
	}

	aiAutoReplyTopic := os.Getenv("TOPIC_AI_AUTOREPLY")
	if aiAutoReplyTopic == "" {
		log.Fatal("TOPIC_AI_AUTOREPLY is not set")
	}

	workerCountStr := os.Getenv("AI_AUTO_REPLY_WORKER_COUNT")
	workerCount := 5 // default
	if workerCountStr != "" {
		if c, err := strconv.Atoi(workerCountStr); err == nil && c > 0 {
			workerCount = c
		}
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

	forwardListener := ws.NewForwardListener(hub, chatService, forwardTopic)
	forwardListener.Start()
	defer forwardListener.Stop()

	inviteListener := ws.NewInviteListener(hub, inviteTopic)
	inviteListener.Start()
	defer inviteListener.Stop()

	aiAutoReplyChan := make(chan kafkaadapter.AiAutoReplyPayload, 100)
	aiAutoReplyConsumer := kafkaadapter.NewAiAutoReplyConsumer(aiAutoReplyTopic)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go aiAutoReplyConsumer.Read(ctx, aiAutoReplyChan)
	defer aiAutoReplyConsumer.Close()

	for i := 0; i < workerCount; i++ {
		go func() {
			for msg := range aiAutoReplyChan {
				ws.HandleAiAutoReply(msg, chatClient, hub)
			}
		}()
	}

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)

	log.Printf("🚀 WebSocket-сервер запущен на :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
