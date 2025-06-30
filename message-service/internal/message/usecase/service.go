package usecase

import (
	"context"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/model"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/ports"
	messagepb "github.com/Vovarama1992/go-ai-messenger/proto/messagepb"
)

type MessageService struct {
	repo       ports.MessageRepo
	userClient ports.UserClient
	chatClient ports.ChatClient
}

func NewMessageService(repo ports.MessageRepo, userClient ports.UserClient) *MessageService {
	return &MessageService{repo: repo, userClient: userClient}
}

func (s *MessageService) SaveMessage(ctx context.Context, msg *model.Message) error {
	msg.CreatedAt = time.Now()
	return s.repo.Save(msg)
}

func (s *MessageService) GetMessagesByChatFiltered(ctx context.Context, chatID int64, limit, offset int) ([]model.Message, error) {
	return s.repo.GetByChat(chatID, limit, offset)
}

func (s *MessageService) GetMessagesByChat(ctx context.Context, chatID int64) ([]*messagepb.ChatMessage, error) {
	const (
		limit  = 10000
		offset = 0
	)

	msgs, err := s.repo.GetByChat(chatID, limit, offset)
	if err != nil {
		return nil, err
	}

	var pbMessages []*messagepb.ChatMessage
	for _, m := range msgs {
		email, err := s.userClient.GetUserEmailByID(ctx, m.SenderID)
		if err != nil {
			return nil, err
		}

		pbMessages = append(pbMessages, &messagepb.ChatMessage{
			SenderId:    m.SenderID,
			SenderEmail: email,
			Text:        m.Text,
			SentAt:      m.CreatedAt.Unix(),
		})
	}

	return pbMessages, nil
}

func (s *MessageService) ResolveThreadInfo(ctx context.Context, threadID string) (*ports.ThreadInfo, error) {
	return s.chatClient.GetThreadInfo(ctx, threadID)
}
