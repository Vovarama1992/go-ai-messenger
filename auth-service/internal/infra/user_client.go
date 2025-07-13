package grpc

import (
	"context"
	"errors"

	"github.com/Vovarama1992/go-ai-messenger/auth-service/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/proto/userpb"
	"github.com/sony/gobreaker"
)

type GrpcUserClient struct {
	client  userpb.UserServiceClient
	breaker *gobreaker.CircuitBreaker
}

func NewGrpcUserClient(client userpb.UserServiceClient, breaker *gobreaker.CircuitBreaker) ports.UserClient {
	return &GrpcUserClient{
		client:  client,
		breaker: breaker,
	}
}

func (g *GrpcUserClient) GetByEmail(ctx context.Context, email string) (int64, string, error) {
	res, err := g.breaker.Execute(func() (interface{}, error) {
		return g.client.GetUserByEmail(ctx, &userpb.GetUserByEmailRequest{Email: email})
	})
	if err != nil {
		return 0, "", err
	}

	userRes, ok := res.(*userpb.GetUserByEmailResponse)
	if !ok {
		return 0, "", errors.New("unexpected response type")
	}

	return userRes.Id, userRes.PasswordHash, nil
}

func (g *GrpcUserClient) Create(ctx context.Context, email, passwordHash string) (int64, error) {
	res, err := g.breaker.Execute(func() (interface{}, error) {
		return g.client.CreateUser(ctx, &userpb.CreateUserRequest{
			Email:        email,
			PasswordHash: passwordHash,
		})
	})
	if err != nil {
		return 0, err
	}

	userRes, ok := res.(*userpb.CreateUserResponse)
	if !ok {
		return 0, errors.New("unexpected response type")
	}

	return userRes.Id, nil
}
