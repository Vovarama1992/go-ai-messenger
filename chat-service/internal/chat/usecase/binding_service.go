package usecase

import (
	"context"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports"
)

type ChatBindingService struct {
	repo ports.ChatBindingRepository
}

func NewChatBindingService(repo ports.ChatBindingRepository) *ChatBindingService {
	return &ChatBindingService{repo: repo}
}

func (s *ChatBindingService) BindUserToChat(ctx context.Context, userID, chatID int64, bindingType model.AIBindingType) error {
	if err := bindingType.IsValid(); err != nil {
		return err
	}

	binding := &model.ChatBinding{
		UserID:    userID,
		ChatID:    chatID,
		Type:      bindingType,
		CreatedAt: time.Now().Unix(),
	}

	return s.repo.Create(ctx, binding)
}

func (s *ChatBindingService) UpdateBinding(ctx context.Context, userID, chatID int64, newType model.AIBindingType) error {
	if err := newType.IsValid(); err != nil {
		return err
	}

	binding := &model.ChatBinding{
		UserID:    userID,
		ChatID:    chatID,
		Type:      newType,
		CreatedAt: time.Now().Unix(),
	}

	return s.repo.Update(ctx, binding)
}

func (s *ChatBindingService) GetBinding(ctx context.Context, userID, chatID int64) (*model.ChatBinding, error) {
	return s.repo.FindByUserAndChat(ctx, userID, chatID)
}

func (s *ChatBindingService) DeleteBinding(ctx context.Context, userID, chatID int64) error {
	return s.repo.Delete(ctx, userID, chatID)
}

func (s *ChatBindingService) GetBindingsByChat(ctx context.Context, chatID int64) ([]*model.ChatBinding, error) {
	return s.repo.FindBindingsByChatID(ctx, chatID)
}
