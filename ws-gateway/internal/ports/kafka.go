package ports

import "context"

type KafkaProducer interface {
	Produce(ctx context.Context, topic string, payload interface{}) error
}
