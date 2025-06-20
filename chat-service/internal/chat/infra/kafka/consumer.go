package kafkaadapter

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(reader *kafka.Reader) *KafkaConsumer {
	return &KafkaConsumer{reader: reader}
}

func (c *KafkaConsumer) StartConsumingThreadResults(ctx context.Context, handler func(model.ThreadResult)) error {
	log.Println("🟢 Kafka consumer started: chat.binding.thread-created")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("❌ Error reading message: %v", err)
			continue
		}

		var result model.ThreadResult
		if err := json.Unmarshal(m.Value, &result); err != nil {
			log.Printf("❌ Invalid JSON in thread-created message: %v", err)
			continue
		}

		log.Printf("📥 Received thread result: chatID=%d userID=%d threadID=%s", result.ChatID, result.UserID, result.ThreadID)
		handler(result)
	}
}
