package kafkaadapter

import (
	"os"

	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter() *kafka.Writer {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	return &kafka.Writer{
		Addr:         kafka.TCP(broker),
		RequiredAcks: kafka.RequireAll,
		Balancer:     &kafka.LeastBytes{},
	}
}

func NewKafkaReader(topic string, groupID string) *kafka.Reader {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{broker},
		Topic:          topic,
		GroupID:        groupID,
		CommitInterval: 0,
	})
}
