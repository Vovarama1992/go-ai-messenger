package grpcadapter

import (
	"context"

	authpb "github.com/Vovarama1992/go-ai-messenger/proto/authpb"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
)

type AuthServiceAdapter struct {
	client authpb.AuthServiceClient
}

func NewAuthService(client authpb.AuthServiceClient) ports.AuthService {
	return &AuthServiceAdapter{client: client}
}

func (a *AuthServiceAdapter) ValidateToken(ctx context.Context, token string) (int64, string, error) {
	resp, err := a.client.ValidateToken(ctx, &authpb.ValidateTokenRequest{Token: token})
	if err != nil {
		return 0, "", err
	}
	return resp.UserId, resp.Email, nil
}
