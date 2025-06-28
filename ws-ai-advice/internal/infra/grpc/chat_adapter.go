package grpcadapter

import (
	"context"
	"fmt"

	chatpb "github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/ports"
)

type ChatServiceAdapter struct {
	client chatpb.ChatServiceClient
}

func NewChatService(client chatpb.ChatServiceClient) ports.ChatService {
	return &ChatServiceAdapter{client: client}
}

func (c *ChatServiceAdapter) GetUserWithChatByThreadID(threadID string) (int64, int64, string, error) {
	resp, err := c.client.GetUserWithChatByThreadID(context.Background(), &chatpb.GetUserWithChatByThreadIDRequest{
		ThreadId: threadID,
	})
	if err != nil {
		return 0, 0, "", err
	}
	if resp == nil {
		return 0, 0, "", fmt.Errorf("nil response from chat service")
	}
	return resp.UserId, resp.ChatId, resp.UserEmail, nil
}
