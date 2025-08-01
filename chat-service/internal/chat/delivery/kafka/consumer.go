package kafkaadapter

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-utils/kafkautil"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader  *kafka.Reader
	retryer *kafkautil.Retry
	breaker *kafkautil.Breaker
	mu      sync.Mutex
}

func NewKafkaConsumer(reader *kafka.Reader) *KafkaConsumer {
	return &KafkaConsumer{
		reader: reader,
	}
}

func (c *KafkaConsumer) WithRetryBreaker(retry *kafkautil.Retry, breaker *kafkautil.Breaker) {
	c.retryer = retry
	c.breaker = breaker
}

func (c *KafkaConsumer) StartConsumingThreadResults(ctx context.Context, handler func(model.ThreadResult) error) error {
	log.Println("🟢 Kafka consumer started: chat.binding.thread-created")

	for {
		var msg kafka.Message
		var err error

		if c.retryer != nil && c.breaker != nil {
			err = c.retryer.Do(func() error {
				return c.breaker.Do(func() error {
					var readErr error
					msg, readErr = c.reader.ReadMessage(ctx)
					return readErr
				})
			})
		} else {
			msg, err = c.reader.ReadMessage(ctx)
		}

		if err != nil {
			log.Printf("❌ Error reading message: %v", err)
			continue
		}

		var result model.ThreadResult
		if err := json.Unmarshal(msg.Value, &result); err != nil {
			log.Printf("❌ Invalid JSON: %v", err)
			continue
		}

		if err := handler(result); err != nil {
			log.Printf("❌ Handler failed: %v", err)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("❌ Commit failed: %v", err)
		}
	}
}
