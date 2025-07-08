package ports

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
)

type GptClient interface {
	CreateThreadForUserAndChat(ctx context.Context, userEmail string, messages []dto.ChatMessage) (string, error)
	SendMessageToThread(ctx context.Context, threadID string, role string, content string) error
	SendMessageAndGetAutoreply(ctx context.Context, threadID string, senderEmail string, text string) (string, error)
	GetAdvice(ctx context.Context, threadID string) (string, error)
}
