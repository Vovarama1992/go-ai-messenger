package ports

import "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/model"

type MessageRepo interface {
	Save(message *model.Message) error
	GetByChat(chatID int64, limit, offset int) ([]model.Message, error)
}
