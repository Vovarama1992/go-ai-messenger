package ports

import "context"

type AuthClient interface {
	ValidateToken(ctx context.Context, token string) (ID int64, Email string, error error)
}
