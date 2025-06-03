package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	chatgrpc "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/grpc"
	chathttp "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/http"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/postgres"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	middleware "github.com/Vovarama1992/go-ai-messenger/pkg/authmiddleware"
	"github.com/Vovarama1992/go-ai-messenger/proto/authpb"
	"github.com/Vovarama1992/go-ai-messenger/proto/chatpb"

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

	// Подключение к базе
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Репозитории и сервисы
	chatRepo := postgres.NewChatRepo(db)
	bindingRepo := postgres.NewChatBindingRepo(db)

	chatService := usecase.NewChatService(chatRepo)
	bindingService := usecase.NewChatBindingService(bindingRepo)

	// Подключение к auth-service по gRPC для middleware
	conn, err := grpc.Dial(authGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to auth-service gRPC: %v", err)
	}
	defer conn.Close()

	authClient := authpb.NewAuthServiceClient(conn)
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	// HTTP сервер
	r := chi.NewRouter()
	r.Use(authMiddleware.Middleware) // применяем auth middleware ко всем роутам

	chathttp.RegisterRoutes(r, chathttp.ChatDeps{
		ChatService:        chatService,
		ChatBindingService: bindingService,
	})

	go func() {
		log.Printf("🚀 HTTP server started on :%s\n", chatHTTPPort)
		if err := http.ListenAndServe(":"+chatHTTPPort, r); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// gRPC сервер
	lis, err := net.Listen("tcp", ":"+chatGRPCPort)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcHandler := chatgrpc.NewChatHandler(chatService, bindingService)
	chatpb.RegisterChatServiceServer(grpcServer, grpcHandler)

	log.Printf("🚀 gRPC server started on :%s\n", chatGRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
