package ports

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) error
	FindByID(ctx context.Context, id int64) (*model.Chat, error)
}
