package kafkaadapter

import (
	"time"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports"
	"github.com/Vovarama1992/go-utils/kafkautil"
	"github.com/segmentio/kafka-go"
)

func NewKafkaProducers(bindingTopic, adviceTopic string, writer *kafka.Writer) (ports.AdvicePublisher, ports.MessageBroker) {
	retry := kafkautil.NewRetry(kafkautil.RetryConfig{
		Attempts: 3,
		Delay:    500 * time.Millisecond,
	})

	breaker := kafkautil.NewBreaker(kafkautil.BreakerConfig{
		MaxFailures: 5,
		Cooldown:    10 * time.Second,
	})

	// Первый продюсер
	chatProducer := NewKafkaProducer(bindingTopic, adviceTopic, writer)
	chatProducer.WithRetryBreaker(retry, breaker)

	// Второй продюсер
	bindingProducer := NewKafkaProducer(bindingTopic, adviceTopic, writer)
	bindingProducer.WithRetryBreaker(retry, breaker)

	return chatProducer, bindingProducer
}
