package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	circuitbreaker "github.com/Vovarama1992/go-ai-messenger/auth-service/internal/config"
	grpchandler "github.com/Vovarama1992/go-ai-messenger/auth-service/internal/delivery/grpc"
	httpadapter "github.com/Vovarama1992/go-ai-messenger/auth-service/internal/delivery/http"
	grpcadapter "github.com/Vovarama1992/go-ai-messenger/auth-service/internal/infra"
	"github.com/Vovarama1992/go-ai-messenger/auth-service/internal/ports"
	auth "github.com/Vovarama1992/go-ai-messenger/auth-service/internal/usecase"
	authpb "github.com/Vovarama1992/go-ai-messenger/proto/authpb"
	"github.com/Vovarama1992/go-ai-messenger/proto/userpb"
	"github.com/Vovarama1992/go-utils/grpcutil"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {

	// Конфиги
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

	// gRPC к user-service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpcadapter.NewUserServiceConn(ctx, userServiceAddr)
	if err != nil {
		log.Fatalf("gRPC подключение к user-service не удалось: %v", err)
	}
	defer conn.Close()

	breaker := circuitbreaker.NewUserServiceBreaker()

	userGrpc := userpb.NewUserServiceClient(conn)
	var userClient ports.UserClient = grpcadapter.NewGrpcUserClient(userGrpc, breaker)

	// Инициализация сервисов и хендлеров
	bcryptLimiter := make(chan struct{}, 4)
	usecase := auth.NewAuthService(userClient, jwtSecret, bcryptLimiter)
	httpHandler := httpadapter.NewHandler(usecase)

	// HTTP маршруты
	mux := http.NewServeMux()
	httpadapter.RegisterRoutes(mux, httpHandler)
	mux.Handle("/docs/", httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
	))
	mux.Handle("/metrics", promhttp.Handler())

	// Запуск HTTP
	go func() {
		log.Printf("Auth HTTP server запущен на :%s\n", httpPort)
		if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
			log.Fatalf("ошибка HTTP-сервера: %v", err)
		}
	}()

	// gRPC сервер
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("не удалось слушать gRPC порт: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcutil.RecoveryInterceptor(),
		),
	)
	grpcHandler := grpchandler.NewHandler(usecase)
	authpb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	log.Printf("auth-service gRPC server запущен на :%s\n", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("ошибка gRPC сервера: %v", err)
	}
}
