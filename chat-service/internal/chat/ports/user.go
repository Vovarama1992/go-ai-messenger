package ports

import "context"

type User struct {
	ID    int64
	Email string
}

type UserClient interface {
	GetUserByID(ctx context.Context, userID int64) (*User, error)
}
