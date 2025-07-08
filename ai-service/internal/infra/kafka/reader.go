package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type Reader struct {
	delegate *kafka.Reader
}

func NewKafkaReader(broker, topic, groupID string) *Reader {
	return &Reader{
		delegate: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     strings.Split(broker, ","),
			Topic:       topic,
			GroupID:     groupID,
			StartOffset: kafka.LastOffset,
			MinBytes:    10e3,
			MaxBytes:    10e6,
			MaxWait:     1 * time.Second,
		}),
	}
}

func (r *Reader) ReadMessage(ctx context.Context) ([]byte, error) {
	msg, err := r.delegate.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return msg.Value, nil
}

func (r *Reader) Close() error {
	return r.delegate.Close()
}
