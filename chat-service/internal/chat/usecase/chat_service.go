package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports"
	"github.com/Vovarama1992/go-utils/ctxutil"
)

type ChatService struct {
	broker          ports.MessageBroker
	chatrepo        ports.ChatRepository
	bindingRepo     ports.ChatBindingRepository
	advicePublisher ports.AdvicePublisher
}

func NewChatService(
	broker ports.MessageBroker,
	chatrepo ports.ChatRepository,
	bindingRepo ports.ChatBindingRepository,
	advicePublisher ports.AdvicePublisher,
) *ChatService {
	return &ChatService{
		broker,
		chatrepo,
		bindingRepo,
		advicePublisher,
	}
}

func (s *ChatService) CreateChat(ctx context.Context, creatorID int64, chatType model.ChatType, memberIDs []int64) (*model.Chat, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, 3)
	defer cancel()

	if err := chatType.IsValid(); err != nil {
		return nil, err
	}

	found := false
	for _, id := range memberIDs {
		if id == creatorID {
			found = true
			break
		}
	}
	if !found {
		memberIDs = append(memberIDs, creatorID)
	}

	chat := &model.Chat{
		CreatorID: creatorID,
		ChatType:  chatType,
		CreatedAt: time.Now().Unix(),
	}

	if err := s.chatrepo.Create(ctx, chat, memberIDs); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) GetChatByID(ctx context.Context, id int64) (*model.Chat, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()

	return s.chatrepo.FindByID(ctx, id)
}

func (s *ChatService) RequestAdvice(ctx context.Context, userID int64, chatID int64) error {
	ctx, cancel := ctxutil.WithTimeout(ctx, 3)
	defer cancel()

	binding, err := s.bindingRepo.FindByUserAndChat(ctx, userID, chatID)
	if err != nil {
		return err
	}
	if binding.BindingType != model.AIBindingAdvice {
		return fmt.Errorf("binding is not of type 'advice'")
	}

	return s.advicePublisher.PublishAdviceRequest(binding.ThreadID)
}

func (s *ChatService) SendInvite(ctx context.Context, chatID int64, userIDs []int64, topic string) error {
	ctx, cancel := ctxutil.WithTimeout(ctx, 3)
	defer cancel()

	if err := s.chatrepo.SendInvite(ctx, chatID, userIDs); err != nil {
		return err
	}

	for _, userID := range userIDs {
		payload := map[string]interface{}{
			"chatId": chatID,
			"userId": userID,
		}
		if err := s.broker.SendInvite(ctx, payload, topic); err != nil {
			log.Printf("failed to send invite kafka message for user %d: %v", userID, err)
		}
	}

	return nil
}

func (s *ChatService) AcceptInvite(ctx context.Context, chatID, userID int64) error {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()

	return s.chatrepo.AcceptInvite(ctx, chatID, userID)
}

func (s *ChatService) GetParticipants(ctx context.Context, chatID int64) ([]int64, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()

	return s.chatrepo.GetChatParticipants(ctx, chatID)
}

func (s *ChatService) GetPendingInvites(ctx context.Context, userID int64) ([]model.Chat, error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, 2)
	defer cancel()

	return s.chatrepo.GetPendingInvites(ctx, userID)
}
