package kafka

import (
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func NewKafkaReader(topic, groupID string) *kafka.Reader {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	log.Printf("📡 Connecting to Kafka: broker=%s topic=%s group=%s", broker, topic, groupID)

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: groupID,
		Topic:   topic,
	})
}
