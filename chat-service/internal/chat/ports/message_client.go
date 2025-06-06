package ports

import (
	"context"

	messagepb "github.com/Vovarama1992/go-ai-messenger/proto/messagepb"
)

type MessageClient interface {
	GetMessagesByChat(ctx context.Context, chatID int64) ([]*messagepb.ChatMessage, error)
}
