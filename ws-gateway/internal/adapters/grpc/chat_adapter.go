package grpcadapter

import (
	"context"

	chatpb "github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
)

type ChatServiceAdapter struct {
	client chatpb.ChatServiceClient
}

func NewChatService(client chatpb.ChatServiceClient) ports.ChatService {
	return &ChatServiceAdapter{client: client}
}

func (a *ChatServiceAdapter) GetBindingsByChat(ctx context.Context, chatID int64) ([]ports.ChatBinding, error) {
	resp, err := a.client.GetBindingsByChat(ctx, &chatpb.GetBindingsByChatRequest{
		ChatId: chatID,
	})
	if err != nil {
		return nil, err
	}

	var bindings []ports.ChatBinding
	for _, b := range resp.Bindings {
		bindings = append(bindings, ports.ChatBinding{
			UserID:      b.UserId,
			BindingType: b.BindingType.String(),
		})
	}

	return bindings, nil
}

func (a *ChatServiceAdapter) GetUsersByChatID(ctx context.Context, chatID int64) ([]int64, error) {
	resp, err := a.client.GetUsersByChatID(ctx, &chatpb.GetUsersByChatIDRequest{
		ChatId: chatID,
	})
	if err != nil {
		return nil, err
	}

	return resp.UserIds, nil
}
