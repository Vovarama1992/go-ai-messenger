package ports

import "context"

type KafkaConsumer interface {
	StartConsuming(ctx context.Context, out chan<- []byte) error
}

type KafkaProducer interface {
	Publish(ctx context.Context, key string, value []byte) error
}
