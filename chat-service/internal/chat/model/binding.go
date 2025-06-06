package model

import "errors"

type AIBindingType string

const (
	AIBindingAdvice    AIBindingType = "advice"
	AIBindingAutoreply AIBindingType = "autoreply"
)

func (a AIBindingType) IsValid() error {
	switch a {
	case AIBindingAdvice, AIBindingAutoreply:
		return nil
	default:
		return errors.New("invalid AI binding type")
	}
}

type ChatBinding struct {
	ChatID    int64
	UserID    int64
	Type      AIBindingType
	CreatedAt int64
}

type AiBindingInitPayload struct {
	ChatID    int64         `json:"chatId"`
	UserID    int64         `json:"userId"`
	Type      string        `json:"type"` // "advice" или "autoreply"
	UserEmail string        `json:"senderEmail"`
	Messages  []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	SenderEmail string `json:"senderEmail"`
	Text        string `json:"text"`
	SentAt      int64  `json:"sentAt"`
}
