package grpc

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports"
	"github.com/Vovarama1992/go-ai-messenger/proto/userpb"
)

type GrpcUserClient struct {
	client userpb.UserServiceClient
}

func NewGrpcUserClient(c userpb.UserServiceClient) ports.UserClient {
	return &GrpcUserClient{client: c}
}

func (g *GrpcUserClient) GetUserByID(ctx context.Context, id int64) (*ports.User, error) {
	resp, err := g.client.GetUserByID(ctx, &userpb.GetUserByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &ports.User{
		ID:    resp.Id,
		Email: resp.Email,
	}, nil
}
