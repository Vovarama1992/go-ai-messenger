package ports

import (
	"context"
)

type ThreadInfo struct {
	UserID    int64
	ChatID    int64
	UserEmail string
}

type ChatClient interface {
	GetThreadInfo(ctx context.Context, threadID string) (*ThreadInfo, error)
}
