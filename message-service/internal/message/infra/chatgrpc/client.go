package chatgrpc

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/ports"
	"github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
)

type ChatClient struct {
	client chatpb.ChatServiceClient
}

func NewChatClient(client chatpb.ChatServiceClient) *ChatClient {
	return &ChatClient{client: client}
}

func (c *ChatClient) GetThreadInfo(ctx context.Context, threadID string) (*ports.ThreadInfo, error) {
	resp, err := c.client.GetUserWithChatByThreadID(ctx, &chatpb.GetUserWithChatByThreadIDRequest{
		ThreadId: threadID,
	})
	if err != nil {
		return nil, err
	}

	return &ports.ThreadInfo{
		UserID:    resp.UserId,
		ChatID:    resp.ChatId,
		UserEmail: resp.UserEmail,
	}, nil
}
