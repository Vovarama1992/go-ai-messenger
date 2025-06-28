package auth

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/ports"
)

type Service struct {
	client ports.AuthClient
}

func NewService(client ports.AuthClient) *Service {
	return &Service{client: client}
}

func (s *Service) ValidateToken(ctx context.Context, token string) (int64, string, error) {
	return s.client.ValidateToken(ctx, token)
}
