package main

import (
	"log"
	"net"
	"os"

	"github.com/Vovarama1992/go-ai-messenger/proto/userpb"
	"github.com/Vovarama1992/go-ai-messenger/user-service/internal/db"
	grpcadapter "github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/adapters/grpc"
	postgresadapter "github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/adapters/postgres"
	"github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/ports"
	userusecase "github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/usecase"
	"google.golang.org/grpc"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL не установлен")
	}

	pool, err := db.NewPool()
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}

	// Адаптер к Postgres
	repo := postgresadapter.NewUserRepository(pool)

	// Usecase
	var service ports.UserService = userusecase.NewUserService(repo)

	// gRPC Handler
	handler := grpcadapter.NewHandler(service)

	grpcPort := os.Getenv("USER_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("не удалось слушать порт: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, handler)

	log.Printf("user-service gRPC server запущен на :%s\n", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("ошибка gRPC сервера: %v", err)
	}
}
