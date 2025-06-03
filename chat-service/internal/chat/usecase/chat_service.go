package usecase

import (
	"context"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports"
)

type ChatService struct {
	repo ports.ChatRepository
}

func NewChatService(repo ports.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) CreateChat(ctx context.Context, userID int64, chatType model.ChatType) (*model.Chat, error) {
	if err := chatType.IsValid(); err != nil {
		return nil, err
	}

	chat := &model.Chat{
		CreatorID: userID,
		Type:      chatType,
		CreatedAt: time.Now().Unix(),
	}

	if err := s.repo.Create(ctx, chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) GetChatByID(ctx context.Context, id int64) (*model.Chat, error) {
	return s.repo.FindByID(ctx, id)
}
