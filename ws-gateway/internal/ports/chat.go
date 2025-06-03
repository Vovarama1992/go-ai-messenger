package ports

import "context"

type ChatService interface {
	GetBindingsByChat(ctx context.Context, chatID int64) ([]ChatBinding, error)
}

type ChatBinding struct {
	UserID int64
	Type   string // "advice" | "autoreply"
}
