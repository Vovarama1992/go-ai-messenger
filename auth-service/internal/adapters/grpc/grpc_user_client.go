package grpc

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/auth-service/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/proto/userpb"
)

type GrpcUserClient struct {
	client userpb.UserServiceClient
}

func NewGrpcUserClient(client userpb.UserServiceClient) ports.UserClient {
	return &GrpcUserClient{client: client}
}

func (g *GrpcUserClient) GetByEmail(ctx context.Context, email string) (int64, string, error) {
	res, err := g.client.GetUserByEmail(ctx, &userpb.GetUserByEmailRequest{Email: email})
	if err != nil {
		return 0, "", err
	}
	return res.Id, res.PasswordHash, nil
}

func (g *GrpcUserClient) Create(ctx context.Context, email, passwordHash string) (int64, error) {
	res, err := g.client.CreateUser(ctx, &userpb.CreateUserRequest{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return 0, err
	}
	return res.Id, nil
}
