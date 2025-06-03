package kafka

import (
	"context"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/dto"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/model"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/usecase"
)

type KafkaMessageProcessor struct {
	service *usecase.MessageService
}

func NewKafkaMessageProcessor(service *usecase.MessageService) *KafkaMessageProcessor {
	return &KafkaMessageProcessor{service: service}
}

func (p *KafkaMessageProcessor) Handle(ctx context.Context, msg dto.IncomingMessage) error {
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
