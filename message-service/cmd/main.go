package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	messagegrpc "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/delivery/grpc"
	deliveryKafka "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/delivery/kafka"
	chatadapter "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/infra/chatgrpc"
	infraKafka "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/infra/postgres"
	useradapter "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/infra/usergrpc"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/streams"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/usecase"
	chatpb "github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
	messagepb "github.com/Vovarama1992/go-ai-messenger/proto/messagepb"
	userpb "github.com/Vovarama1992/go-ai-messenger/proto/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	topicPersist := os.Getenv("TOPIC_MESSAGE_PERSIST")

	msgChan := streams.MessageChan

	consumerCount, err := strconv.Atoi(os.Getenv("TOPIC_MESSAGE_PERSIST_CONSUMER_COUNT"))
	if err != nil || consumerCount <= 0 {
		log.Printf("⚠️  Invalid or unset TOPIC_MESSAGE_PERSIST_CONSUMER_COUNT, defaulting to 1")
		consumerCount = 1
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	repo := postgres.NewMessageRepo(db)

	userServiceAddr := os.Getenv("USER_GRPC_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = "localhost:50052"
	}
	conn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to user-service: %v", err)
	}
	defer conn.Close()

	userGRPC := userpb.NewUserServiceClient(conn)
	userClient := useradapter.NewUserClient(userGRPC)

	chatServiceAddr := os.Getenv("CHAT_GRPC_ADDR")
	if chatServiceAddr == "" {
		chatServiceAddr = "localhost:50053"
	}
	chatConn, err := grpc.Dial(chatServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to chat-service: %v", err)
	}
	defer chatConn.Close()

	chatGRPC := chatpb.NewChatServiceClient(chatConn)
	chatClient := chatadapter.NewChatClient(chatGRPC)

	messageService := usecase.NewMessageService(repo, userClient)
	processor := postgres.NewDefaultMessageProcessor(messageService, chatClient)

	workerCount := 4
	if wcStr := os.Getenv("CHAT_MESSAGE_PERSIST_WORKER_COUNT"); wcStr != "" {
		if wcParsed, err := strconv.Atoi(wcStr); err == nil && wcParsed > 0 {
			workerCount = wcParsed
		} else {
			log.Printf("⚠️  Invalid CHAT_MESSAGE_PERSIST_WORKER_COUNT, using default (%d)", workerCount)
		}
	}

	for i := 0; i < consumerCount; i++ {
		reader := infraKafka.NewKafkaReader(topicPersist, "message-group")
		deliveryKafka.StartKafkaConsumer(ctx, reader, msgChan)
	}

	deliveryKafka.StartMessageWorkers(ctx, wg, processor, msgChan, workerCount)

	grpcServer := grpc.NewServer()
	messageHandler := messagegrpc.NewMessageHandler(messageService)
	messagepb.RegisterMessageServiceServer(grpcServer, messageHandler)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("🛑 Shutting down...")
	cancel()
	grpcServer.GracefulStop()
	wg.Wait()
	log.Println("✅ Shutdown complete")
}
