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
	infraKafka "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/infra/postgres"
	useradapter "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/infra/usergrpc"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/usecase"
	messagepb "github.com/Vovarama1992/go-ai-messenger/proto/messagepb"
	userpb "github.com/Vovarama1992/go-ai-messenger/proto/userpb"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

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

	userConn, err := grpc.Dial(userServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Failed to connect to user-service: %v", err)
	}
	defer userConn.Close()

	userGRPC := userpb.NewUserServiceClient(userConn) // как у тебя уже есть
	userClient := useradapter.NewUserClient(userGRPC)

	messageService := usecase.NewMessageService(repo, userClient)
	processor := deliveryKafka.NewKafkaMessageProcessor(messageService)

	workerCount := 4
	if wcStr := os.Getenv("CHAT_MESSAGE_PERSIST_WORKER_COUNT"); wcStr != "" {
		if wcParsed, err := strconv.Atoi(wcStr); err == nil && wcParsed > 0 {
			workerCount = wcParsed
		} else {
			log.Printf("⚠️  Invalid CHAT_MESSAGE_PERSIST_WORKER_COUNT, using default (%d)", workerCount)
		}
	}

	reader := infraKafka.NewKafkaReader("chat.message.persist", "message-group")
	infraKafka.StartKafkaConsumer(ctx, reader)
	deliveryKafka.StartMessageWorkers(ctx, wg, processor, workerCount)

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
