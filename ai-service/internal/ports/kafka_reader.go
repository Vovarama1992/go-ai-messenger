package ports

import "context"

type KafkaReader interface {
	ReadMessage(ctx context.Context) ([]byte, error)
	Close() error
}
