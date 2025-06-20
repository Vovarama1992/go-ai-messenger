package usecase

import (
	"context"
	"log"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports"
)

type ChatBindingService struct {
	repo          ports.ChatBindingRepository
	broker        ports.MessageBroker
	messageClient ports.MessageClient
}

func NewChatBindingService(repo ports.ChatBindingRepository, broker ports.MessageBroker, client ports.MessageClient) *ChatBindingService {
	return &ChatBindingService{repo, broker, client}
}

func (s *ChatBindingService) UpdateBinding(ctx context.Context, userEmail string, userID, chatID int64, newType model.AIBindingType) error {
	if err := newType.IsValid(); err != nil {
		return err
	}

	binding := &model.ChatBinding{
		UserID:    userID,
		ChatID:    chatID,
		Type:      newType,
		CreatedAt: time.Now().Unix(),
	}

	_, err := s.repo.FindByUserAndChat(ctx, userID, chatID)
	if err != nil {
		// Привязки нет — создаём и пушим ивент с историей
		if err := s.repo.Create(ctx, binding); err != nil {
			return err
		}

		messagesPb, err := s.messageClient.GetMessagesByChat(ctx, chatID)
		if err != nil {
			return err
		}

		// Конвертируем protobuf-сообщения в model.ChatMessage
		messages := make([]model.ChatMessage, len(messagesPb))
		for i, m := range messagesPb {
			messages[i] = model.ChatMessage{
				SenderEmail: m.SenderEmail,
				Text:        m.Text,
				SentAt:      m.SentAt,
			}
		}

		payload := model.AiBindingInitPayload{
			ChatID:    chatID,
			UserID:    userID,
			Type:      string(newType),
			UserEmail: userEmail,
			Messages:  messages,
		}

		go s.broker.SendAiBindingInit(context.Background(), payload)
		return nil
	}

	// Привязка уже есть — обновляем
	return s.repo.Update(ctx, binding)
}

func (s *ChatBindingService) HandleThreadCreated(res model.ThreadResult) {
	ctx := context.Background()

	err := s.repo.UpdateThreadID(ctx, res.ChatID, res.UserID, res.ThreadID)
	if err != nil {
		// тут можно кастомную обработку: retry, log-level, метрики и т.д.
		log.Printf("❌ Failed to update threadID: %v", err)
	} else {
		log.Printf("✅ ThreadID updated: chatID=%d, userID=%d, threadID=%s", res.ChatID, res.UserID, res.ThreadID)
	}
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
