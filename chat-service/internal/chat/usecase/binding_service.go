package usecase

import (
	"context"
	"log"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports"
	"github.com/Vovarama1992/go-utils/ctxutil"
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
	ctx, cancel := ctxutil.WithTimeout(ctx, 3)
	defer cancel()
	if err := newType.IsValid(); err != nil {
		return err
	}

	binding := &model.ChatBinding{
		UserID:      userID,
		ChatID:      chatID,
		BindingType: newType,
		CreatedAt:   time.Now().Unix(),
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
				Text:        m.Content,
				SentAt:      m.SentAt,
			}
		}

		payload := model.AiBindingInitPayload{
			ChatID:      chatID,
			UserID:      userID,
			BindingType: string(newType),
			UserEmail:   userEmail,
			Messages:    messages,
		}

		go s.broker.SendAiBindingInit(context.Background(), payload)
		return nil
	}

	// Привязка уже есть — обновляем
	return s.repo.Update(ctx, binding)
}

func (s *ChatBindingService) HandleThreadCreated(res model.ThreadResult) error {
	ctx, cancel := ctxutil.WithTimeout(context.Background(), 3)
	defer cancel()

	err := s.repo.UpdateThreadID(ctx, res.ChatID, res.UserID, res.ThreadID)
	if err != nil {
		log.Printf("❌ Failed to update threadID: %v", err)
		return err
	}

	log.Printf("✅ ThreadID updated: chatID=%d, userID=%d, threadID=%s", res.ChatID, res.UserID, res.ThreadID)
	return nil
}
func (s *ChatBindingService) GetBinding(ctx context.Context, userID, chatID int64) (*model.ChatBinding, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()
	return s.repo.FindByUserAndChat(ctx, userID, chatID)
}

func (s *ChatBindingService) DeleteBinding(ctx context.Context, userID, chatID int64) error {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()
	return s.repo.Delete(ctx, userID, chatID)
}

func (s *ChatBindingService) GetBindingsByChat(ctx context.Context, chatID int64) ([]*model.ChatBinding, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()
	return s.repo.FindBindingsByChatID(ctx, chatID)
}

func (s *ChatBindingService) FindByThreadID(ctx context.Context, threadID string) (*model.ChatBinding, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()
	return s.repo.FindByThreadID(ctx, threadID)
}
