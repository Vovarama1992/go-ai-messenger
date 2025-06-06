package ports

import "context"

type UserClient interface {
	GetUserEmailByID(ctx context.Context, id int64) (string, error)
}
