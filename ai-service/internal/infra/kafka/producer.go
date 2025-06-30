package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topics []string
}

func NewKafkaWriter(broker string, topics ...string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Balancer: &kafka.LeastBytes{},
		},
		topics: topics,
	}
}

func (p *Producer) Publish(ctx context.Context, key string, value []byte) error {
	var msgs []kafka.Message
	for _, topic := range p.topics {
		msgs = append(msgs, kafka.Message{
			Topic: topic,
			Key:   []byte(key),
			Value: value,
		})
	}
	return p.writer.WriteMessages(ctx, msgs...)
}
