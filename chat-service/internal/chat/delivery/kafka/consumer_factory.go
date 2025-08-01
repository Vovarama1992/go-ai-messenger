package kafkaadapter

import (
	"time"

	"github.com/Vovarama1992/go-utils/kafkautil"
	"github.com/segmentio/kafka-go"
)

func NewThreadCreatedConsumer(reader *kafka.Reader) *KafkaConsumer {
	consumer := NewKafkaConsumer(reader)

	retry := kafkautil.NewRetry(kafkautil.RetryConfig{
		Attempts: 3,
		Delay:    1 * time.Second,
	})

	breaker := kafkautil.NewBreaker(kafkautil.BreakerConfig{
		MaxFailures: 5,
		Cooldown:    10 * time.Second,
	})

	consumer.WithRetryBreaker(retry, breaker)
	return consumer
}
