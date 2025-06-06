package ports

import "context"

type AuthService interface {
	ValidateToken(ctx context.Context, token string) (int64, string, error)
}
