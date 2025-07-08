package ports

import "context"

type KafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...[]byte) error
}
