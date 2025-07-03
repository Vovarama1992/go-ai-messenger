package ports

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat, memberIDs []int64) error
	FindByID(ctx context.Context, id int64) (*model.Chat, error)
	SendInvite(ctx context.Context, chatID int64, userIDs []int64) error
	AcceptInvite(ctx context.Context, chatID, userID int64) error
	GetChatParticipants(ctx context.Context, chatID int64) ([]int64, error)
	GetPendingInvites(ctx context.Context, userID int64) ([]model.Chat, error)
}
