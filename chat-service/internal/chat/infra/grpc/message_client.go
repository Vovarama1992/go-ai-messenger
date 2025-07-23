package grpc

import (
	"context"

	messagepb "github.com/Vovarama1992/go-ai-messenger/proto/messagepb"
)

type GrpcMessageClient struct {
	client messagepb.MessageServiceClient
}

func NewGrpcMessageClient(client messagepb.MessageServiceClient) *GrpcMessageClient {
	return &GrpcMessageClient{client: client}
}

func (g *GrpcMessageClient) GetMessagesByChat(ctx context.Context, chatID int64) ([]*messagepb.ChatMessage, error) {
	resp, err := g.client.GetMessagesByChat(ctx, &messagepb.GetMessagesRequest{ChatId: chatID})
	if err != nil {
		return nil, err
	}
	return resp.Messages, nil
}
