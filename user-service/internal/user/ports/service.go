package ports

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/user-service/internal/user/model"
)

type UserService interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, email, passwordHash string) (*model.User, error)
}
