package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	chatgrpc "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/grpc"
	grpcclient "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/grpc"
	chathttp "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/http"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/postgres"
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	middleware "github.com/Vovarama1992/go-ai-messenger/pkg/authmiddleware"
	authpb "github.com/Vovarama1992/go-ai-messenger/proto/authpb"
	chatpb "github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
	messagepb "github.com/Vovarama1992/go-ai-messenger/proto/messagepb"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	chatHTTPPort := os.Getenv("CHAT_HTTP_PORT")
	if chatHTTPPort == "" {
		chatHTTPPort = "8081"
	}

	chatGRPCPort := os.Getenv("CHAT_GRPC_PORT")
	if chatGRPCPort == "" {
		chatGRPCPort = "50053"
	}

	authGRPCAddr := os.Getenv("AUTH_SERVICE_GRPC_ADDR")
	if authGRPCAddr == "" {
		authGRPCAddr = "auth-service:50052"
	}

	messageGRPCAddr := os.Getenv("MESSAGE_SERVICE_GRPC_ADDR")
	if messageGRPCAddr == "" {
		log.Fatal("MESSAGE_SERVICE_GRPC_ADDR is not set")
	}

	topic := os.Getenv("TOPIC_AI_BINDING_INIT")
	if topic == "" {
		log.Fatal("TOPIC_AI_BINDING_INIT is not set")
	}

	broker := kafkaadapter.NewKafkaProducer(topic)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	chatRepo := postgres.NewChatRepo(db)
	bindingRepo := postgres.NewChatBindingRepo(db)

	msgConn, err := grpc.Dial(messageGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to message-service gRPC: %v", err)
	}
	defer msgConn.Close()

	msgClient := grpcclient.NewGrpcMessageClient(messagepb.NewMessageServiceClient(msgConn))

	chatService := usecase.NewChatService(chatRepo)
	bindingService := usecase.NewChatBindingService(bindingRepo, broker, msgClient)

	conn, err := grpc.Dial(authGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to auth-service gRPC: %v", err)
	}
	defer conn.Close()

	authClient := authpb.NewAuthServiceClient(conn)
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	r := chi.NewRouter()
	r.Use(authMiddleware.Middleware)

	chathttp.RegisterRoutes(r, chathttp.ChatDeps{
		ChatService:        chatService,
		ChatBindingService: bindingService,
	})

	go func() {
		log.Printf("\U0001F680 HTTP server started on :%s\n", chatHTTPPort)
		if err := http.ListenAndServe(":"+chatHTTPPort, r); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+chatGRPCPort)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcHandler := chatgrpc.NewChatHandler(chatService, bindingService)
	chatpb.RegisterChatServiceServer(grpcServer, grpcHandler)

	log.Printf("\U0001F680 gRPC server started on :%s\n", chatGRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
