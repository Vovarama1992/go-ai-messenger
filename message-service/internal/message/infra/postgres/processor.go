package postgres

import (
	"context"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/dto"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/model"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/ports"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/usecase"
)

type DefaultMessageProcessor struct {
	service     *usecase.MessageService
	chatService ports.ChatClient
}

func NewDefaultMessageProcessor(service *usecase.MessageService, chat ports.ChatClient) *DefaultMessageProcessor {
	return &DefaultMessageProcessor{
		service:     service,
		chatService: chat,
	}
}
func (p *DefaultMessageProcessor) Handle(ctx context.Context, msg dto.IncomingMessage) error {
	// Обогащение, если не хватает данных
	if msg.SenderID == 0 && msg.ChatID == 0 && msg.ThreadID != "" {
		info, err := p.chatService.GetThreadInfo(ctx, msg.ThreadID)
		if err != nil {
			log.Printf("❌ Failed to fetch thread info: %v", err)
			return err
		}
		msg.SenderID = info.UserID
		msg.ChatID = info.ChatID
		msg.AIGenerated = true
	}

	message := &model.Message{
		ChatID:      msg.ChatID,
		SenderID:    msg.SenderID,
		Text:        msg.Text,
		AIGenerated: msg.AIGenerated,
	}

	if err := p.service.SaveMessage(ctx, message); err != nil {
		log.Printf("❌ DB save error: %v", err)
		return err
	}

	log.Printf("✅ Saved message ID %d: %+v", message.ID, message)
	return nil
}
