package ports

import "context"

type ChatService interface {
	GetBindingsByChat(ctx context.Context, chatID int64) ([]ChatBinding, error)
	GetUsersByChatID(ctx context.Context, chatID int64) ([]int64, error) //
}

type ChatBinding struct {
	UserID      int64
	ChatID      int64
	ThreadID    string
	BindingType string // "advice" | "autoreply"
}
