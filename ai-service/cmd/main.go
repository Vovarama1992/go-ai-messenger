package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	advice "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/app/advice"
	binding "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/app/binding"
	feed "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/app/feed"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/infra/gpt"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/infra/kafka"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleShutdown(cancel)

	// ENV
	broker := os.Getenv("KAFKA_BROKER")
	apiKey := os.Getenv("OPENAI_API_KEY")
	assistantID := os.Getenv("OPENAI_ASSISTANT_ID")

	topicBindingInit := os.Getenv("TOPIC_AI_BINDING_INIT")
	topicThreadCreated := os.Getenv("TOPIC_AI_THREAD_CREATED")
	topicAiFeed := os.Getenv("TOPIC_AI_FEED")
	topicAiAutoreply := os.Getenv("TOPIC_AI_AUTOREPLY")
	topicPersist := os.Getenv("TOPIC_MESSAGE_PERSIST")
	topicAdviceReq := os.Getenv("TOPIC_AI_ADVICE_REQUEST")
	topicAdviceResp := os.Getenv("TOPIC_AI_ADVICE_RESPONSE")

	if topicAiAutoreply == "" {
		log.Fatal("❌ TOPIC_AI_AUTOREPLY is not set")
	}

	// GPT client
	gptClient := gpt.NewClient(apiKey, assistantID)

	// Kafka readers
	bindingReader := kafka.NewKafkaReader(broker, topicBindingInit, "ai-binding-init-group")
	feedReader := kafka.NewKafkaReader(broker, topicAiFeed, "ai-feed-group")
	adviceReader := kafka.NewKafkaReader(broker, topicAdviceReq, "ai-advice-group")

	// Kafka writers
	threadWriter := kafka.NewKafkaWriter(broker, topicThreadCreated)
	autoreplyWriter := kafka.NewKafkaWriter(broker, topicAiAutoreply, topicPersist)
	adviceWriter := kafka.NewKafkaWriter(broker, topicAdviceResp)

	// Binding pipeline
	binding.RunAiCreateBindingReaderFromKafka(ctx, 10, bindingReader)
	binding.RunAiCreateBindingReaderFromGpt(ctx, 10, gptClient)
	binding.RunAiCreateBindingWriterToKafka(ctx, threadWriter, topicThreadCreated)

	// Feed pipeline
	feed.RunAiFeedReaderFromKafka(ctx, 10, feedReader)
	feed.RunAiFeedReaderFromGpt(ctx, 10, gptClient)
	feed.RunAiFeedWriterToKafka(ctx, autoreplyWriter)

	// Advice pipeline
	advice.RunAdviceReaderFromKafka(ctx, 10, adviceReader)
	advice.RunAdviceReaderFromGpt(ctx, 10, gptClient)
	advice.RunAdviceWriterToKafka(ctx, adviceWriter, topicAdviceResp)

	<-ctx.Done()
	log.Println("🛑 Main stopped")
}

func handleShutdown(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Println("📴 Shutdown signal received")
	cancel()
}
