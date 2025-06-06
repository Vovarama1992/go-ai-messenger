package kafka

import (
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func NewKafkaReader(broker, topic string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     strings.Split(broker, ","),
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.LastOffset,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,
	})
}

func NewKafkaWriter(broker, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(strings.Split(broker, ",")...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}
