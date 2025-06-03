package ports

import (
	"context"
)

type UserClient interface {
	GetByEmail(ctx context.Context, email string) (id int64, passwordHash string, err error)
	Create(ctx context.Context, email, passwordHash string) (id int64, err error)
}
