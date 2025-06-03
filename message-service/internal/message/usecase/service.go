package usecase

import (
	"context"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/model"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/ports"
)

type MessageService struct {
	repo ports.MessageRepo
}

func NewMessageService(repo ports.MessageRepo) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) SaveMessage(ctx context.Context, msg *model.Message) error {
	msg.CreatedAt = time.Now()
	return s.repo.Save(msg)
}
