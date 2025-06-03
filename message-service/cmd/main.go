package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	deliveryKafka "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/delivery/kafka"
	infraKafka "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/infra/postgres"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/usecase"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Connect to Postgres
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Compose layers
	repo := postgres.NewMessageRepo(db)
	service := usecase.NewMessageService(repo)
	processor := deliveryKafka.NewKafkaMessageProcessor(service)

	// Parse worker count from env
	workerCount := 4 // default
	if wcStr := os.Getenv("CHAT_MESSAGE_PERSIST_WORKER_COUNT"); wcStr != "" {
		if wcParsed, err := strconv.Atoi(wcStr); err == nil && wcParsed > 0 {
			workerCount = wcParsed
		} else {
			log.Printf("⚠️  Invalid CHAT_MESSAGE_PERSIST_WORKER_COUNT, using default (%d)", workerCount)
		}
	}

	// Start Kafka consumer and worker pool
	reader := infraKafka.NewKafkaReader("chat.message.persist", "message-group")
	infraKafka.StartKafkaConsumer(ctx, reader)
	deliveryKafka.StartMessageWorkers(ctx, wg, processor, workerCount)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("🛑 Shutting down...")
	cancel()
	wg.Wait()
	log.Println("✅ Shutdown complete")
}
