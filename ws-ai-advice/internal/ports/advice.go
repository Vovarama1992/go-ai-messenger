package ports

import (
	"context"
)

type AdviceConsumer interface {
	Start(context.Context) error
	Close() error
}
