package usecase

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports"
	"github.com/Vovarama1992/go-utils/ctxutil"
)

type UserService struct {
	client ports.UserClient
}

func NewUserService(client ports.UserClient) *UserService {
	return &UserService{client: client}
}

func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*ports.User, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()
	return s.client.GetUserByID(ctx, userID)
}
