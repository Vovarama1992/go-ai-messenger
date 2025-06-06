package model

type ChatMessage struct {
	SenderEmail string `json:"sender_email"`
	Text        string `json:"text"`
	SentAt      int64  `json:"sent_at"`
}

type AiBindingInitPayload struct {
	ChatID    int64         `json:"chat_id"`
	UserID    int64         `json:"user_id"`
	UserEmail string        `json:"user_email"`
	Type      string        `json:"type"`
	Messages  []ChatMessage `json:"messages"`
}
