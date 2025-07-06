package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	chatgrpc "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/grpc"
	chathttp "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/http"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/postgres"
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	middleware "github.com/Vovarama1992/go-ai-messenger/pkg/authmiddleware"
	authpb "github.com/Vovarama1992/go-ai-messenger/proto/authpb"
	chatpb "github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
	messagepb "github.com/Vovarama1992/go-ai-messenger/proto/messagepb"
	userpb "github.com/Vovarama1992/go-ai-messenger/proto/userpb"

	_ "github.com/Vovarama1992/go-ai-messenger/chat-service/docs"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	// ENV
	dbURL := os.Getenv("DATABASE_URL")
	chatHTTPPort := getEnvOrDefault("CHAT_HTTP_PORT", "8081")
	chatGRPCPort := getEnvOrDefault("CHAT_GRPC_PORT", "50053")
	authGRPCAddr := getEnvOrDefault("AUTH_SERVICE_GRPC_ADDR", "auth-service:50052")
	messageGRPCAddr := os.Getenv("MESSAGE_SERVICE_GRPC_ADDR")
	userGRPCAddr := os.Getenv("USER_SERVICE_GRPC_ADDR")
	bindingTopic := os.Getenv("TOPIC_AI_BINDING_INIT")
	resultTopic := os.Getenv("TOPIC_AI_THREAD_CREATED")
	adviceTopic := os.Getenv("TOPIC_AI_ADVICE_REQUEST")

	inviteTopic := os.Getenv("TOPIC_CHAT_INVITE")
	if inviteTopic == "" {
		log.Fatal("❌ TOPIC_CHAT_INVITE env is not set")
	}

	if dbURL == "" || bindingTopic == "" || messageGRPCAddr == "" || resultTopic == "" {
		log.Fatal("❌ One or more required environment variables are not set")
	}

	// DB
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ failed to connect to DB: %v", err)
	}
	defer db.Close()

	var chatRepo ports.ChatRepository = postgres.NewChatRepo(db)
	var bindingRepo ports.ChatBindingRepository = postgres.NewChatBindingRepo(db)

	// Kafka producer
	writer := kafkaadapter.NewKafkaWriter()
	var chatproducer ports.AdvicePublisher = kafkaadapter.NewKafkaProducer(bindingTopic, adviceTopic, writer)
	var bindingproducer ports.MessageBroker = kafkaadapter.NewKafkaProducer(bindingTopic, adviceTopic, writer)

	// Kafka consumer (чтение chat.binding.thread-created)
	reader := kafkaadapter.NewKafkaReader(resultTopic, "chat-service")
	consumer := kafkaadapter.NewKafkaConsumer(reader)

	// Message gRPC client
	msgConn, err := grpc.Dial(messageGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ failed to connect to message-service gRPC: %v", err)
	}
	userConn, err := grpc.Dial(userGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ failed to connect to user-service gRPC: %v", err)
	}
	defer msgConn.Close()
	defer userConn.Close()
	msgClient := chatgrpc.NewGrpcMessageClient(messagepb.NewMessageServiceClient(msgConn))
	userClient := chatgrpc.NewGrpcUserClient(userpb.NewUserServiceClient(userConn))
	userService := usecase.NewUserService(userClient)

	// Services
	chatService := usecase.NewChatService(bindingproducer, chatRepo, bindingRepo, chatproducer)
	bindingService := usecase.NewChatBindingService(bindingRepo, bindingproducer, msgClient)

	go consumer.StartConsumingThreadResults(ctx, bindingService.HandleThreadCreated)

	// Auth middleware
	conn, err := grpc.Dial(authGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ failed to connect to auth-service gRPC: %v", err)
	}
	defer conn.Close()
	authClient := authpb.NewAuthServiceClient(conn)
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	// HTTP API
	r := chi.NewRouter()
	r.Use(authMiddleware.Middleware)
	chathttp.RegisterRoutes(r, chathttp.ChatDeps{
		ChatService:        chatService,
		ChatBindingService: bindingService,
	}, inviteTopic)

	r.Get("/swagger/*", httpSwagger.Handler())

	go func() {
		log.Printf("🚀 HTTP server started on :%s", chatHTTPPort)
		if err := http.ListenAndServe(":"+chatHTTPPort, r); err != nil {
			log.Fatalf("❌ HTTP server error: %v", err)
		}
	}()

	// gRPC API
	lis, err := net.Listen("tcp", ":"+chatGRPCPort)
	if err != nil {
		log.Fatalf("❌ failed to listen on gRPC port: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcHandler := chatgrpc.NewChatHandler(chatService, bindingService, userService)
	chatpb.RegisterChatServiceServer(grpcServer, grpcHandler)

	log.Printf("🚀 gRPC server started on :%s", chatGRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ gRPC server error: %v", err)
	}
}

func getEnvOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
