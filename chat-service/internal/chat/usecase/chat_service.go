package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports"
)

type ChatService struct {
	chatrepo        ports.ChatRepository
	bindingRepo     ports.ChatBindingRepository
	advicePublisher ports.AdvicePublisher
}

func NewChatService(
	chatrepo ports.ChatRepository,
	bindingRepo ports.ChatBindingRepository,
	advicePublisher ports.AdvicePublisher,
) *ChatService {
	return &ChatService{
		chatrepo:        chatrepo,
		bindingRepo:     bindingRepo,
		advicePublisher: advicePublisher,
	}
}

func (s *ChatService) CreateChat(ctx context.Context, userID int64, chatType model.ChatType) (*model.Chat, error) {
	if err := chatType.IsValid(); err != nil {
		return nil, err
	}

	chat := &model.Chat{
		CreatorID: userID,
		Type:      chatType,
		CreatedAt: time.Now().Unix(),
	}

	if err := s.chatrepo.Create(ctx, chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) GetChatByID(ctx context.Context, id int64) (*model.Chat, error) {
	return s.chatrepo.FindByID(ctx, id)
}

func (s *ChatService) RequestAdvice(ctx context.Context, userID int64, chatID int64) error {
	binding, err := s.bindingRepo.FindByUserAndChat(ctx, userID, chatID)
	if err != nil {
		return err
	}
	if binding.Type != model.AIBindingAdvice {
		return fmt.Errorf("binding is not of type 'advice'")
	}

	return s.advicePublisher.PublishAdviceRequest(binding.ThreadID)
}
