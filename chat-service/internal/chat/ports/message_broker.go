package ports

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
)

type MessageBroker interface {
	SendAiBindingInit(ctx context.Context, payload model.AiBindingInitPayload) error
}
