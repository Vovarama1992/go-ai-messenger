package ports

import "context"

type AuthService interface {
	ValidateToken(ctx context.Context, token string) (int64, string, error)
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, email, password string) (int64, error)
}
