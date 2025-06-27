package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"

	socketio "github.com/googollee/go-socket.io"

	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/stream"
)

func main() {
	port := os.Getenv("WS_AI_ADVICE_PORT")
	if port == "" {
		log.Fatal("WS_AI_ADVICE_PORT is not set")
	}

	// WebSocket
	server := socketio.NewServer(nil)
	server.OnConnect("/", func(c socketio.Conn) error {
		log.Println("✅ connected:", c.ID())
		return nil
	})
	server.OnDisconnect("/", func(c socketio.Conn, reason string) {
		log.Println("❌ disconnected:", c.ID(), reason)
	})

	// Kafka Reader
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092"
	}
	topic := "chat.message.ai-advice"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokers},
		GroupID: "ws-ai-advice",
		Topic:   topic,
	})

	consumer := kafkaadapter.NewAdviceConsumer(reader)

	// Контекст + graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Fatalf("❌ consumer start failed: %v", err)
		}
	}()

	// Подписка на канал → пушим в сокет
	go func() {
		for advice := range stream.PendingAdviceChan {
			// TODO: заменить userID или threadID на socket room или mapping
			server.BroadcastToNamespace("/", "gpt-advice", advice)
		}
	}()

	// HTTP server
	go func() {
		http.Handle("/socket.io/", server)
		log.Printf("🚀 ws-ai-advice up on :%s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("❌ ListenAndServe: %v", err)
		}
	}()

	// Ждём сигнал завершения
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("🛑 shutdown initiated")
	consumer.Close()
	server.Close()
}
