package chat

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/ports"
)

type Service struct {
	client ports.ChatService
}

func NewService(client ports.ChatService) *Service {
	return &Service{client: client}
}

func (s *Service) GetUserWithChatByThreadID(ctx context.Context, threadID string) (int64, int64, string, error) {
	return s.client.GetUserWithChatByThreadID(threadID)
}
