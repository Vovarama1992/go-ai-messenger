package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	grpcadapter "github.com/Vovarama1992/go-ai-messenger/auth-service/internal/auth/adapters/grpc"
	httpadapter "github.com/Vovarama1992/go-ai-messenger/auth-service/internal/auth/adapters/http"
	"github.com/Vovarama1992/go-ai-messenger/auth-service/internal/auth/ports"
	auth "github.com/Vovarama1992/go-ai-messenger/auth-service/internal/auth/usecase"
	authpb "github.com/Vovarama1992/go-ai-messenger/proto/authpb"
	"github.com/Vovarama1992/go-ai-messenger/proto/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET не установлен")
	}

	userServiceAddr := os.Getenv("USER_SERVICE_GRPC_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = "localhost:50051"
	}

	httpPort := os.Getenv("AUTH_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	grpcPort := os.Getenv("AUTH_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		userServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("gRPC подключение к user-service не удалось: %v", err)
	}
	defer conn.Close()

	userGrpc := userpb.NewUserServiceClient(conn)
	var userClient ports.UserClient = grpcadapter.NewGrpcUserClient(userGrpc)

	usecase := auth.NewAuthService(userClient, jwtSecret)

	httpHandler := httpadapter.NewHandler(usecase)
	http.HandleFunc("/login", httpHandler.Login)
	http.HandleFunc("/register", httpHandler.Register)

	go func() {
		log.Printf("Auth HTTP server запущен на :%s\n", httpPort)
		if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
			log.Fatalf("ошибка HTTP-сервера: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("не удалось слушать gRPC порт: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcHandler := grpcadapter.NewHandler(usecase)
	authpb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	log.Printf("auth-service gRPC server запущен на :%s\n", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("ошибка gRPC сервера: %v", err)
	}
}
