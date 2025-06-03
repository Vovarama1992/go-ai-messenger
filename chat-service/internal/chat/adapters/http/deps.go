package http

import (
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
)

type ChatDeps struct {
	ChatService        *usecase.ChatService
	ChatBindingService *usecase.ChatBindingService
}
