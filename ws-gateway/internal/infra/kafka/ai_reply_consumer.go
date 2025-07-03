package kafkaadapter

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type AiAutoReplyPayload struct {
	ThreadID string `json:"threadId"`
	Text     string `json:"text"`
}

type AiAutoReplyConsumer struct {
	reader *kafka.Reader
}

func NewAiAutoReplyConsumer(topic string) *AiAutoReplyConsumer {
	broker := os.Getenv("KAFKA_BROKERS")
	if broker == "" {
		broker = "kafka:9092"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: "ai-auto-reply-group",
		Topic:   topic,
	})

	return &AiAutoReplyConsumer{reader: reader}
}

func (c *AiAutoReplyConsumer) Read(ctx context.Context, out chan<- AiAutoReplyPayload) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("❌ Kafka read error: %v", err)
			continue
		}

		var msg AiAutoReplyPayload
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("❌ Unmarshal error: %v", err)
			continue
		}

		out <- msg
	}
}

func (c *AiAutoReplyConsumer) Close() error {
	return c.reader.Close()
}
