package ports

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
)

type ChatBindingRepository interface {
	Create(ctx context.Context, binding *model.ChatBinding) error
	UpdateThreadID(ctx context.Context, chatID, userID int64, threadID string) error
	FindByUserAndChat(ctx context.Context, userID, chatID int64) (*model.ChatBinding, error)
	Update(ctx context.Context, binding *model.ChatBinding) error
	Delete(ctx context.Context, userID, chatID int64) error
	FindBindingsByChatID(ctx context.Context, chatID int64) ([]*model.ChatBinding, error)
	FindByThreadID(ctx context.Context, threadID string) (*model.ChatBinding, error)
}
