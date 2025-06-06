package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	infraKafka "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/infra/kafka"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	go handleShutdown(cancel)

	// Init Kafka reader
	topic := os.Getenv("TOPIC_AI_BINDING_INIT")
	if topic == "" {
		log.Fatal("❌ TOPIC_AI_BINDING_INIT not set")
	}

	reader := infraKafka.NewKafkaReader(topic, "ai-binding-init-group")
	infraKafka.StartKafkaConsumer(ctx, reader)

	// Process messages
	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Service stopped")
			return

		case msg := <-infraKafka.MessageChan:
			log.Printf("📨 Received AiBindingInitPayload: chatID=%d, userID=%d, messages=%d\n",
				msg.ChatID, msg.UserID, len(msg.Messages))
		}
	}
}

func handleShutdown(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Println("📴 Caught shutdown signal")
	cancel()
}
